# 1. log in
curl -X POST http://127.0.0.1:8080/users/yirzhou/login
# 2. log out
curl -X POST http://127.0.0.1:8080/users/yirzhou/logout
# 3. expecting error because user not logged in
curl -X PUT "http://127.0.0.1:8080/users/yirzhou/topics/add?topic=amazon"
# 4. log in user and follow a topic
curl -X POST http://127.0.0.1:8080/users/yirzhou/login
curl -X PUT "http://127.0.0.1:8080/users/yirzhou/topics/add?topic=amazon"
# 5. follow topics
curl -X PUT "http://127.0.0.1:8080/users/yirzhou/topics/remove?topic=amazon"
curl -X PUT "http://127.0.0.1:8080/users/yirzhou/topics/remove?topic=alibaba"
curl -X POST http://127.0.0.1:8080/users/yirzhou/logout
# 6. Get New Posts after logging in
curl -X POST http://127.0.0.1:8080/users/yirzhou/login
curl -X GET http://127.0.0.1:8080/users/yirzhou/new_posts
curl -X POST http://127.0.0.1:8080/users/yirzhou/logout
curl -X POST http://127.0.0.1:8080/users/yirzhou/login
curl -X GET http://127.0.0.1:8080/users/yirzhou/new_posts