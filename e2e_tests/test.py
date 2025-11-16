import pytest
import requests
import os


BASE = os.getenv("HOST_ADDRESS")

# -----------------------------
# Fixtures for unique teams/users
# -----------------------------
@pytest.fixture(scope="session")
def pr_team_data():
    team_name = "pr-team-unique"
    members = [
        {"user_id": "p1", "username": "user_p1", "is_active": True},
        {"user_id": "p2", "username": "user_p2", "is_active": True},
        {"user_id": "p3", "username": "user_p3", "is_active": True},
    ]
    r = requests.post(f"{BASE}/team/add", json={"team_name": team_name, "members": members})
    assert r.status_code in (201, 400)  # already exists is fine
    return {"team_name": team_name, "members": members}


# -----------------------------
# Helpers
# -----------------------------
def create_pr(pr_id, name, author_id):
    return requests.post(
        f"{BASE}/pullRequest/create",
        json={"pull_request_id": pr_id, "pull_request_name": name, "author_id": author_id},
    )


def merge_pr(pr_id):
    return requests.post(f"{BASE}/pullRequest/merge", json={"pull_request_id": pr_id})


def reassign_pr(pr_id, old_reviewer_id):
    return requests.post(f"{BASE}/pullRequest/reassign", json={
        "pull_request_id": pr_id,
        "old_reviewer_id": old_reviewer_id
    })


# -----------------------------
# PR TESTS
# -----------------------------
def test_pr_create_and_conflict(pr_team_data):
    # create PR
    r1 = create_pr("pr-unique-1", "PR1", "p1")
    assert r1.status_code == 201

    # create duplicate PR
    r2 = create_pr("pr-unique-1", "PR1-DUP", "p1")
    assert r2.status_code == 409
    assert r2.json()["error"]["code"] == "PR_EXISTS"


def test_pr_create_author_not_found():
    r = create_pr("pr-unique-x", "PRX", "ghost")
    assert r.status_code == 404
    assert r.json()["error"]["code"] == "NOT_FOUND"


def test_pr_merge(pr_team_data):
    create_pr("pr-unique-merge", "MergeTest", "p1")
    r = merge_pr("pr-unique-merge")
    assert r.status_code == 200
    assert r.json()["pr"]["status"] == "MERGED"


def test_pr_merge_not_found():
    r = merge_pr("pr-ghost")
    assert r.status_code == 404


def test_pr_reassign_ok(pr_team_data):
    r = create_pr("pr-unique-reassign", "ReassignTest", "p1")
    assert r.status_code == 201
    assigned = r.json()["pr"]["assigned_reviewers"]

    old = assigned[0]
    r = reassign_pr("pr-unique-reassign", old)
    assert r.status_code == 200
    assert r.json()["replaced_by"] != old


def test_pr_reassign_not_assigned(pr_team_data):
    create_pr("pr-na", "PRNA", "p1")
    r = reassign_pr("pr-na", "p1")  # author is not reviewer
    assert r.status_code == 409
    assert r.json()["error"]["code"] == "NOT_ASSIGNED"


def test_pr_reassign_no_candidate():
    # create a team with only one active user
    team_name = "nc-team"
    members = [
        {"user_id": "nc1", "username": "nc_user1", "is_active": True},
        {"user_id": "nc2", "username": "nc_user2", "is_active": False},
    ]
    requests.post(f"{BASE}/team/add", json={"team_name": team_name, "members": members})

    r = create_pr("pr-nc", "PRNC", "nc2")
    assigned = r.json()["pr"]["assigned_reviewers"]
    old = assigned[0]

    r = reassign_pr("pr-nc", old)
    assert r.status_code == 409
    assert r.json()["error"]["code"] == "NO_CANDIDATE"


def test_pr_reassign_on_merged(pr_team_data):
    create_pr("pr-merged", "PRMerged", "p1")
    merge_pr("pr-merged")

    # pick a reviewer
    pr_info = requests.get(f"{BASE}/users/getReview", params={"user_id": "p2"}).json()
    if pr_info["pull_requests"]:
        old_reviewer = pr_info["pull_requests"][0]["pull_request_id"]
    else:
        old_reviewer = "p2"

    r = reassign_pr("pr-merged", "p2")
    assert r.status_code == 409
    assert r.json()["error"]["code"] == "PR_MERGED"
