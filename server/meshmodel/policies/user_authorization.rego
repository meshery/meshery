package user_authorization

default allow = false

allow {
    input.role_names[_] == "team admin"
}
