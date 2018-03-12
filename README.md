# auth microservice of recreationroom project
# Use cases
## 注册
用户输入用户名、密码、昵称、邮箱，若用户名、邮箱未被占用，则注册成功。

## 销号
用户登录后，可以删除自己的账号。

# REST API
GET /users              UserService.List        List all users
POST /users             UserService.Register    User registration
GET /users/<id>         UserService.Show        Show user
PATCH /users/<id>       UserService.Edit        Edit user
DELETE /users/<id>      UserService.Delete      Unregister

GET /sessions/<id>      SessionService.Get      Load session from redis; check if user is logged in
POST /sessions          SessionService.Create   Create session; User login
DELETE /sessions/<id>   SessionService.Delete   Delete session; User logout
