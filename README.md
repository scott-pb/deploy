# deploy
```yaml
ip: 192.168.0.99
port: 8088
# 账号
accounts:
  - username: admin
    password: admin123
  - username: scott
    password: scott123
  - username: stan
    password: stan1996
  - username: bob
    password: bob123

# 环境
admin_test:
  project_config:
    project_path: project/admin_test
    project_name: soga_admin
    bin_path: project/admin_test/bin
    git_url: http://192.168.0.13/platform/soga_admin.git
    # ssh
  client_config:
    host: 192.168.0.30
    port: 22
    user: root
    password: admin@123
    # git账号密码
  git_config:
    user_name: bin.peng
    pass_word: pb123654
  # 打包相关
  build_configs:
    - env: admin
      mod_path: project/admin_test/soga_admin
      bin_name: bin/soga_admin
      name: soga_admin
  zip_file_path: bin/admin_server.zip
  zip_name: admin_server.zip
  server_path: /root/soga/soga_im_admin_server

admin_release:
  project_config:
    project_path: project/admin_release
    bin_path: project/admin_release/soga_admin/bin
    project_name: soga_admin
    git_url: http://192.168.0.13/platform/soga_admin.git
  client_config:
    host: 47.243.81.177
    port: 22
    user: root
    password: X@135791667
  git_config:
    user_name: bin.peng
    pass_word: pb123654
  build_configs:
    - env: admin
      mod_path: project/admin_release/soga_admin
      bin_name: bin/soga_admin
      name: soga_admin
  zip_file_path: bin/admin_server.zip
  zip_name: admin_server.zip
  server_path: /root/soga_admin_server

im_enterprise_test:
  zip_file_path: bin/im_enterprise.zip
  zip_name: im_enterprise.zip
  server_path: /root/soga/soga_im_enterprise/bin
  project_config:
    project_path: project/im_enterprise_test
    project_name: soga_im_enterprise
    bin_path: project/im_enterprise_test/soga_im_enterprise/bin
    git_url: http://192.168.0.13/server/soga_im_enterprise.git
  client_config:
    host: 192.168.0.30
    port: 22
    user: root
    password: admin@123
  git_config:
    user_name: bin.peng
    pass_word: pb123654
  build_configs:
    - env: soga_api_chat
      mod_path: project/im_enterprise_test/soga_im_enterprise/cmd/api/chat
      bin_name: bin/soga_api_chat
      name: soga_api_chat
    - env: soga_api_chatroom
      mod_path: project/im_enterprise_test/soga_im_enterprise/cmd/api/chatroom
      bin_name: bin/soga_api_chatroom
      name: soga_api_chatroom
    - env: soga_rpc_chat
      mod_path: project/im_enterprise_test/soga_im_enterprise/cmd/rpc/chat
      bin_name: bin/soga_rpc_chat
      name: soga_rpc_chat
    - env: soga_rpc_game
      mod_path: project/im_enterprise_test/soga_im_enterprise/cmd/rpc/game
      bin_name: bin/soga_rpc_game
      name: soga_rpc_game
    - env: soga_cron
      mod_path: project/im_enterprise_test/soga_im_enterprise/cmd/cron
      bin_name: bin/soga_cron
      bin: soga_cron
      name: soga_cron
    - env: soga_tool
      mod_path: project/im_enterprise_test/soga_im_enterprise/cmd/tool
      bin_name: bin/soga_tool
      bin: soga_tool
      name: soga_tool

im_enterprise_release:
  zip_file_path: bin/im_enterprise.zip
  zip_name: im_enterprise.zip
  server_path: /root/soga/soga_im_enterprise/bin
  project_config:
    project_path: project/im_enterprise_release
    bin_path: project/im_enterprise_release/soga_im_enterprise/bin
    project_name: soga_im_enterprise
    git_url: http://192.168.0.13/server/soga_im_enterprise.git
  client_config:
    host: 47.243.81.177
    port: 22
    user: root
    password: X@135791667
  git_config:
    user_name: bin.peng
    pass_word: pb123654
  build_configs:
    - env: soga_api_chat
      mod_path: project/im_enterprise_release/soga_im_enterprise/cmd/api/chat
      bin_name: bin/soga_api_chat
      name: soga_api_chat
    - env: soga_api_chatroom
      mod_path: project/im_enterprise_release/soga_im_enterprise/cmd/api/chatroom
      bin_name: bin/soga_api_chatroom
      name: soga_api_chatroom
    - env: soga_rpc_chat
      mod_path: project/im_enterprise_release/soga_im_enterprise/cmd/rpc/chat
      bin_name: bin/soga_rpc_chat
      name: soga_rpc_chat
    - env: soga_rpc_game
      mod_path: project/im_enterprise_release/soga_im_enterprise/cmd/rpc/game
      bin_name: bin/soga_rpc_game
      name: soga_rpc_game
    - env: soga_cron
      mod_path: project/im_enterprise_release/soga_im_enterprise/cmd/cron
      bin_name: bin/soga_cron
      bin: soga_cron
      name: soga_cron
    - env: soga_tool
      mod_path: project/im_enterprise_release/soga_im_enterprise/cmd/tool
      bin_name: bin/soga_tool
      bin: soga_tool
      name: soga_tool

im_server_test:
  zip_file_path: bin/im_server.zip
  zip_name: im_server.zip
  server_path: /root/soga/soga_im_server/bin
  project_config:
    project_path: project/im_server_test
    bin_path: project/im_server_test/soga_im_server/bin
    project_name: soga_im_server
    git_url: http://192.168.0.13/server/soga_im_server.git
  client_config:
    host: 192.168.0.30
    port: 22
    user: root
    password: admin@123
  git_config:
    user_name: bin.peng
    pass_word: pb123654
  build_configs:
    - env: soga_im_api
      mod_path: project/im_server_test/soga_im_server/cmd/api/soga_im_api
      bin_name: bin/soga_im_api
      name: soga_im_api
    - env: soga_im_msg_gateway
      mod_path: project/im_server_test/soga_im_server/cmd/api/soga_im_msg_gateway
      bin_name: bin/soga_im_msg_gateway
      name: soga_im_msg_gateway
    - env: soga_im_msg_transfer
      mod_path: project/im_server_test/soga_im_server/cmd/api/soga_im_msg_transfer
      bin_name: bin/soga_im_msg_transfer
      name: soga_im_msg_transfer
    - env: soga_im_push
      mod_path: project/im_server_test/soga_im_server/cmd/api/soga_im_push
      bin_name: bin/soga_im_push
      name: soga_im_push
    - env: soga_im_rpc_auth
      mod_path: project/im_server_test/soga_im_server/cmd/rpc/soga_im_rpc_auth
      bin_name: bin/soga_im_rpc_auth
      name: soga_im_rpc_auth
    - env: soga_im_rpc_cache
      mod_path: project/im_server_test/soga_im_server/cmd/rpc/soga_im_rpc_cache
      bin_name: bin/soga_im_rpc_cache
      name: soga_im_rpc_cache
    - env: soga_im_rpc_conversation
      mod_path: project/im_server_test/soga_im_server/cmd/rpc/soga_im_rpc_conversation
      bin_name: bin/soga_im_rpc_conversation
      name: soga_im_rpc_conversation
    - env: soga_im_rpc_friend
      mod_path: project/im_server_test/soga_im_server/cmd/rpc/soga_im_rpc_friend
      bin_name: bin/soga_im_rpc_friend
      name: soga_im_rpc_friend
    - env: soga_im_rpc_group
      mod_path: project/im_server_test/soga_im_server/cmd/rpc/soga_im_rpc_group
      bin_name: bin/soga_im_rpc_group
      name: soga_im_rpc_group
    - env: soga_im_rpc_msg
      mod_path: project/im_server_test/soga_im_server/cmd/rpc/soga_im_rpc_msg
      bin_name: bin/soga_im_rpc_msg
      name: soga_im_rpc_msg
    - env: soga_im_rpc_office
      mod_path: project/im_server_test/soga_im_server/cmd/rpc/soga_im_rpc_office
      bin_name: bin/soga_im_rpc_office
      name: soga_im_rpc_office
    - env: soga_im_rpc_organization
      mod_path: project/im_server_test/soga_im_server/cmd/rpc/soga_im_rpc_organization
      bin_name: bin/soga_im_rpc_organization
      name: soga_im_rpc_organization
    - env: soga_im_rpc_user
      mod_path: project/im_server_test/soga_im_server/cmd/rpc/soga_im_rpc_user
      bin_name: bin/soga_im_rpc_user
      name: soga_im_rpc_user

im_server_release:
  zip_file_path: bin/im_server.zip
  zip_name: im_server.zip
  server_path: /root/soga/soga_im_server/bin
  project_config:
    project_path: project/im_server_release
    bin_path: project/im_server_release/soga_im_server/bin
    project_name: soga_im_server
    git_url: http://192.168.0.13/server/soga_im_server.git
  client_config:
    host: 47.243.81.177
    port: 22
    user: root
    password: X@135791667
  git_config:
    user_name: bin.peng
    pass_word: pb123654
  build_configs:
    - env: soga_im_api
      mod_path: project/im_server_release/soga_im_server/cmd/api/soga_im_api
      bin_name: bin/soga_im_api
      name: soga_im_api
    - env: soga_im_msg_gateway
      mod_path: project/im_server_release/soga_im_server/cmd/api/soga_im_msg_gateway
      bin_name: bin/soga_im_msg_gateway
      name: soga_im_msg_gateway
    - env: soga_im_msg_transfer
      mod_path: project/im_server_release/soga_im_server/cmd/api/soga_im_msg_transfer
      bin_name: bin/soga_im_msg_transfer
      name: soga_im_msg_transfer
    - env: soga_im_push
      mod_path: project/im_server_release/soga_im_server/cmd/api/soga_im_push
      bin_name: bin/soga_im_push
      name: soga_im_push
    - env: soga_im_rpc_auth
      mod_path: project/im_server_release/soga_im_server/cmd/rpc/soga_im_rpc_auth
      bin_name: bin/soga_im_rpc_auth
      name: soga_im_rpc_auth
    - env: soga_im_rpc_cache
      mod_path: project/im_server_release/soga_im_server/cmd/rpc/soga_im_rpc_cache
      bin_name: bin/soga_im_rpc_cache
      name: soga_im_rpc_cache
    - env: soga_im_rpc_conversation
      mod_path: project/im_server_release/soga_im_server/cmd/rpc/soga_im_rpc_conversation
      bin_name: bin/soga_im_rpc_conversation
      name: soga_im_rpc_conversation
    - env: soga_im_rpc_friend
      mod_path: project/im_server_release/soga_im_server/cmd/rpc/soga_im_rpc_friend
      bin_name: bin/soga_im_rpc_friend
      name: soga_im_rpc_friend
    - env: soga_im_rpc_group
      mod_path: project/im_server_release/soga_im_server/cmd/rpc/soga_im_rpc_group
      bin_name: bin/soga_im_rpc_group
      name: soga_im_rpc_group
    - env: soga_im_rpc_msg
      mod_path: project/im_server_release/soga_im_server/cmd/rpc/soga_im_rpc_msg
      bin_name: bin/soga_im_rpc_msg
      name: soga_im_rpc_msg
    - env: soga_im_rpc_office
      mod_path: project/im_server_release/soga_im_server/cmd/rpc/soga_im_rpc_office
      bin_name: bin/soga_im_rpc_office
      name: soga_im_rpc_office
    - env: soga_im_rpc_organization
      mod_path: project/im_server_release/soga_im_server/cmd/rpc/soga_im_rpc_organization
      bin_name: bin/soga_im_rpc_organization
      name: soga_im_rpc_organization
    - env: soga_im_rpc_user
      mod_path: project/im_server_release/soga_im_server/cmd/rpc/soga_im_rpc_user
      bin_name: bin/soga_im_rpc_user
      name: soga_im_rpc_user

admin_production:
  project_config:
    project_path: project/admin_production
    bin_path: project/admin_production/soga_admin/bin
    project_name: soga_admin
    git_url: http://192.168.0.13/platform/soga_admin.git
  client_config:
    host: 47.243.81.177
    port: 22
    user: root
    password: X@135791667
  git_config:
    user_name: bin.peng
    pass_word: pb123654
  build_configs:
    - env: admin
      mod_path: project/admin_production/soga_admin
      bin_name: bin/soga_admin
      name: soga_admin
  zip_file_path: bin/admin_server.zip
  zip_name: admin_server.zip
  server_path: /root/soga/soga_im_admin_server

im_enterprise_production:
  zip_file_path: bin/im_enterprise.zip
  zip_name: im_enterprise.zip
  server_path: /root/soga/soga_im_enterprise/bin
  project_config:
    project_path: project/im_enterprise_production
    bin_path: project/im_enterprise_production/soga_im_enterprise/bin
    project_name: soga_im_enterprise
    git_url: http://192.168.0.13/server/soga_im_enterprise.git
  client_config:
    host: 47.243.81.177
    port: 22
    user: root
    password: X@135791667
  git_config:
    user_name: bin.peng
    pass_word: pb123654
  build_configs:
    - env: soga_api_chat
      mod_path: project/im_enterprise_production/soga_im_enterprise/cmd/api/chat
      bin_name: bin/soga_api_chat
      name: soga_api_chat
    - env: soga_api_chatroom
      mod_path: project/im_enterprise_production/soga_im_enterprise/cmd/api/chatroom
      bin_name: bin/soga_api_chatroom
      name: soga_api_chatroom
    - env: soga_rpc_chat
      mod_path: project/im_enterprise_production/soga_im_enterprise/cmd/rpc/chat
      bin_name: bin/soga_rpc_chat
      name: soga_rpc_chat
    - env: soga_rpc_game
      mod_path: project/im_enterprise_production/soga_im_enterprise/cmd/rpc/game
      bin_name: bin/soga_rpc_game
      name: soga_rpc_game
    - env: soga_cron
      mod_path: project/im_enterprise_production/soga_im_enterprise/cmd/cron
      bin_name: bin/soga_cron
      bin: soga_cron
      name: soga_cron
    - env: soga_tool
      mod_path: project/im_enterprise_production/soga_im_enterprise/cmd/tool
      bin_name: bin/soga_tool
      bin: soga_tool
      name: soga_tool

im_server_production:
  zip_file_path: bin/im_server.zip
  zip_name: im_server.zip
  server_path: /root/soga/soga_im_server/bin
  project_config:
    project_path: project/im_server_production
    bin_path: project/im_server_production/soga_im_server/bin
    project_name: soga_im_server
    git_url: http://192.168.0.13/server/soga_im_server.git
  client_config:
    host: 47.243.81.177
    port: 22
    user: root
    password: X@135791667
  git_config:
    user_name: bin.peng
    pass_word: pb123654
  build_configs:
    - env: soga_im_api
      mod_path: project/im_server_production/soga_im_server/cmd/api/soga_im_api
      bin_name: bin/soga_im_api
      name: soga_im_api
    - env: soga_im_msg_gateway
      mod_path: project/im_server_production/soga_im_server/cmd/api/soga_im_msg_gateway
      bin_name: bin/soga_im_msg_gateway
      name: soga_im_msg_gateway
    - env: soga_im_msg_transfer
      mod_path: project/im_server_production/soga_im_server/cmd/api/soga_im_msg_transfer
      bin_name: bin/soga_im_msg_transfer
      name: soga_im_msg_transfer
    - env: soga_im_push
      mod_path: project/im_server_production/soga_im_server/cmd/api/soga_im_push
      bin_name: bin/soga_im_push
      name: soga_im_push
    - env: soga_im_rpc_auth
      mod_path: project/im_server_production/soga_im_server/cmd/rpc/soga_im_rpc_auth
      bin_name: bin/soga_im_rpc_auth
      name: soga_im_rpc_auth
    - env: soga_im_rpc_cache
      mod_path: project/im_server_production/soga_im_server/cmd/rpc/soga_im_rpc_cache
      bin_name: bin/soga_im_rpc_cache
      name: soga_im_rpc_cache
    - env: soga_im_rpc_conversation
      mod_path: project/im_server_production/soga_im_server/cmd/rpc/soga_im_rpc_conversation
      bin_name: bin/soga_im_rpc_conversation
      name: soga_im_rpc_conversation
    - env: soga_im_rpc_friend
      mod_path: project/im_server_production/soga_im_server/cmd/rpc/soga_im_rpc_friend
      bin_name: bin/soga_im_rpc_friend
      name: soga_im_rpc_friend
    - env: soga_im_rpc_group
      mod_path: project/im_server_production/soga_im_server/cmd/rpc/soga_im_rpc_group
      bin_name: bin/soga_im_rpc_group
      name: soga_im_rpc_group
    - env: soga_im_rpc_msg
      mod_path: project/im_server_production/soga_im_server/cmd/rpc/soga_im_rpc_msg
      bin_name: bin/soga_im_rpc_msg
      name: soga_im_rpc_msg
    - env: soga_im_rpc_office
      mod_path: project/im_server_production/soga_im_server/cmd/rpc/soga_im_rpc_office
      bin_name: bin/soga_im_rpc_office
      name: soga_im_rpc_office
    - env: soga_im_rpc_organization
      mod_path: project/im_server_production/soga_im_server/cmd/rpc/soga_im_rpc_organization
      bin_name: bin/soga_im_rpc_organization
      name: soga_im_rpc_organization
    - env: soga_im_rpc_user
      mod_path: project/im_server_production/soga_im_server/cmd/rpc/soga_im_rpc_user
      bin_name: bin/soga_im_rpc_user
      name: soga_im_rpc_user

admin_ui:
  project_config:
    project_path: project/admin_ui
    project_name: soga_admin_ui
    git_url: http://192.168.0.13/platform/soga_admin_ui.git
  test_client_config:
    host: 192.168.0.30
    port: 22
    user: root
    password: admin@123
  release_client_config:
    host: 47.243.81.177
    port: 22
    user: root
    password: X@135791667
  git_config:
    user_name: bin.peng
    pass_word: pb123654
  zip_name: dist.zip
  test_server_path: /root/soga/soga_im_web
  release_server_path: /root/soga_admin_web
```