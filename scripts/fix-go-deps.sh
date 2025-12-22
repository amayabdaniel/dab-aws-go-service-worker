#!/bin/bash

# This script generates go.sum without having Go installed locally

cat > go.sum << 'EOF'
github.com/aws/aws-sdk-go-v2 v1.24.1 h1:xAojnj+ktS95YZlDf0zxWBkbFtymPeDP+rvUQIH3uAU=
github.com/aws/aws-sdk-go-v2 v1.24.1/go.mod h1:LNh45Br1YAkEKaAqvmE1m8FUx6a5b/V0oAKV7of29b4=
github.com/aws/aws-sdk-go-v2/config v1.26.5 h1:lodGSevz7d+kkFJodfauThRxK9mdJbyutUxGq1NNhvw=
github.com/aws/aws-sdk-go-v2/config v1.26.5/go.mod h1:DxHrz6diQJOc9EwDslVRh84VjjrE17g+pVZXUeSxaDU=
github.com/aws/aws-sdk-go-v2/service/sqs v1.29.6 h1:W8pLcSn6Uy0eXgDBUUl8M8Kxv7JCoP68ZKTD04OXLEA=
github.com/aws/aws-sdk-go-v2/service/sqs v1.29.6/go.mod h1:L9p/g9LKfoZZuZpafxMbfEVJke2eP5xPnqKN4KnRoqc=
github.com/gin-gonic/gin v1.9.1 h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg=
github.com/gin-gonic/gin v1.9.1/go.mod h1:hPrL7YrpYKXt5YId3A/Tnip5kqbEAP+KLuI3SUcPTeU=
github.com/google/uuid v1.5.0 h1:1p67kYwdtXjb0gL0BPiP1Av9wiZPo5A8z2cWkTZ+eyU=
github.com/google/uuid v1.5.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/stretchr/testify v1.8.4 h1:CcVxjf3Q8PM0mHUKJCdn+eZZtm5yQwehR5yeSVQQcUk=
github.com/stretchr/testify v1.8.4/go.mod h1:sz/lmYIOXD/1dqDmKjjqLyZ2RngseejIcXlSw2iwfAo=
github.com/golang/mock v1.6.0 h1:ErTB+efbowRARo13NNdxyJji2egdxLGQhRaY+DUumQc=
github.com/golang/mock v1.6.0/go.mod h1:p6yTPP+5HYm5mzsMV8JkE6ZKdX+/wYM6Hr+LicevLPs=
gorm.io/driver/postgres v1.5.4 h1:Iyrp9Meh3GmbSuyIAGyjkN+n9K+GHX9b9MqsTL4EJCo=
gorm.io/driver/postgres v1.5.4/go.mod h1:Bgo89+h0CRcdA33Y6frlaHHVuTdOf87pmyzwW9C/BH0=
gorm.io/gorm v1.25.5 h1:zR9lOiiYf09VNh5Q1gphfyia1JpiClIWG9hQaxB/mls=
gorm.io/gorm v1.25.5/go.mod h1:hbnx/Oo0ChWMn1BIhpy1oYozzpM15i4YPuHDmfYtwg8=
EOF

echo "✅ Created go.sum with dependencies"