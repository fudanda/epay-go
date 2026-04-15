# EPay Go

一个基于 Go + Gin + PostgreSQL + Redis + Vue 3 的支付系统示例项目，提供管理后台、商户中心、统一下单、通道管理，以及订单 / 退款 / 结算流程。

cd web && npm run build
go build -o epay-server cmd/server/main.go


## 技术栈

- 后端：Go、Gin、GORM、PostgreSQL、Redis
- 前端：Vue 3、Vite、Arco Design
- 部署：Docker Compose、Nginx / Caddy

## 目录说明

- `cmd/server`：服务启动入口
- `internal`：后端核心业务
- `web`：前端代码
- `deploy/nginx`：Nginx 配置示例
- `deploy/caddy`：Caddy 配置
- `docker-compose.yml`：默认部署
- `docker-compose.prod.caddy.yml`：Caddy 直接对外的生产部署示例

## 快速开始

### 1. 准备环境变量

```bash
cp .env.example .env
```

然后按需修改数据库、Redis、JWT、默认管理员和支付渠道配置。

### 2. 启动项目

```bash
docker compose up -d --build
```

默认包含以下服务：

- `postgres`
- `redis`
- `backend`
- `frontend`

默认端口：

- `80`：前端
- `8080`：后端
- `5432`：PostgreSQL
- `6379`：Redis

### 常用访问入口

部署完成后，可直接访问以下前端路径：

- **管理员登录**：`/admin/login`
- **商户注册**：`/merchant/register`
- **商户登录**：`/merchant/login`

## 环境变量

参考 `.env.example`。常用变量包括：

- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `REDIS_PASSWORD`
- `JWT_SECRET`
- `DEFAULT_ADMIN_USERNAME`
- `DEFAULT_ADMIN_PASSWORD`
- `SITE_ADDRESS`
- `ACME_EMAIL`
- `ALIPAY_APP_ID`
- `ALIPAY_PRIVATE_KEY`
- `ALIPAY_PUBLIC_KEY`
- `WECHAT_APP_ID`
- `WECHAT_MCH_ID`
- `WECHAT_API_KEY`

系统首次启动且数据库中没有管理员时，会使用 `DEFAULT_ADMIN_USERNAME` 和 `DEFAULT_ADMIN_PASSWORD` 初始化默认管理员。

## 部署

### Caddy 直接对外

适用于宿主机未占用 `80/443`：

```bash
docker compose -f docker-compose.prod.caddy.yml up -d --build
```

### 宿主机 Nginx 反向代理

适用于宿主机已有 Nginx：

```bash
docker compose up -d --build
```

然后由宿主机 Nginx 反代到容器端口，示例配置见：

- `deploy/nginx/host.prod.conf.example`

## 支付参数说明

支付通道和支付场景是分开的：

- `type` / `pay_type`：决定渠道，例如 `wxpay`、`alipay`
- `pay_method`：决定场景，例如 `native`、`scan`、`h5`、`jsapi`、`web`

### 关键规则

- `type=native` **不允许单独使用**，因为无法判断是微信还是支付宝
- 后端支持显式别名，并会自动归一化：
  - `WX_NATIVE`
  - `WX_JSAPI`
  - `WX_H5`
  - `ALIPAY_SCAN`
  - `ALIPAY_H5`
  - `ALIPAY_WEB`

### 推荐传法

- **微信 Native**
  - `type=WX_NATIVE`
  - 或 `type=wxpay&pay_method=native`

- **微信 JSAPI**
  - `type=WX_JSAPI`
  - 或 `type=wxpay&pay_method=jsapi`

- **支付宝扫码**
  - `type=ALIPAY_SCAN`
  - 或 `type=alipay&pay_method=scan`

- **支付宝 H5**
  - `type=ALIPAY_H5`
  - 或 `type=alipay&pay_method=h5`

- **支付宝网页支付**
  - `type=ALIPAY_WEB`
  - 或 `type=alipay&pay_method=web`

## 构建说明

如果在中国大陆网络环境构建，可以在 `.env` 中设置：

```env
GOPROXY=https://goproxy.cn,direct
```


