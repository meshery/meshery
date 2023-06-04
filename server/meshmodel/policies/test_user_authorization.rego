package test_user_authorization

import data.user_authorization.allow

test_if_user_allowed {
    allow with input as {
        "id": "349b7693-7fee",
        "user_id": "suhailkna@gmail.com",
        "first_name": "Suhail",
        "last_name": "Khan",
        "avatar_url": "https://lh3.googleusercontent.com/a/ALm5wu3dIgTAEhSkqsGQ-wcfQDBhK2tT2eWBeWjT2xih2w=s96-c",
        "provider": "meshery-cloud",
        "email": "suhailkna@gmail.com",
        "role_names": [
            "meshmap",
            "team admin"
        ]
    }
}

test_if_user_not_allowed{
    allow with input as {"role_names": ["admin", "manager"]}
}