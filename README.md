# auth microservice of recreationroom project
# Use cases
## 注册
用户输入用户名、密码、昵称、邮箱，若用户名、邮箱未被占用，则注册成功。

## 销号
用户登录后，可以删除自己的账号。

# REST API
GET /users              user.Service.List               List all users
POST /users             user.Service.Create             User registration
GET /users/<id>         user.Service.Get                Get user
PATCH /users/<id>       user.Service.Update             Update user
DELETE /users/<id>      user.Service.Delete             Unregister

POST /users/<id>/password user.Service.CreatePassword   Reset password(need password reset token)
PUT /users/<id>/password user.Service.UpdatePassword    Update user's password(need current password)
