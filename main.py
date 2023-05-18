import streamlit as st
import requests
import json
import os

# 应用程序存储文件夹
APP_FOLDER = "var"


# 获取应用程序列表
def get_app_list():
    app_list = []
    for file in os.listdir(APP_FOLDER):
        if file.endswith(".json"):
            with open(os.path.join(APP_FOLDER, file), "r") as f:
                app_data = json.load(f)
                app_list.append(app_data)
    return app_list


# 创建应用程序
def create_app(app_name, doip, run_script):
    app_data = {"name": app_name, "doip": doip, "run_script": run_script}
    with open(os.path.join(APP_FOLDER, f"{app_name}.json"), "w") as f:
        json.dump(app_data, f)


# 运行应用程序
def run_app(app_name):
    # 在这里执行HTTP请求以启动应用程序
    body = {"doip": "abacus"}
    headers = {"Content-Type": "application/json"}
    response = requests.post(
        f"http://localhost:8080/run", headers=headers, data=json.dumps(body)
    )
    return response.status_code


# 主应用程序界面
st.title("科研数据复现平台")

# 侧边栏
st.sidebar.title("视图选择")
view_options = ["查看应用程序", "创建应用程序"]
selected_view = st.sidebar.selectbox("选择一个视图：", view_options)

# 查看应用程序视图
if selected_view == "查看应用程序":
    app_list = get_app_list()
    app_selected = st.selectbox("选择一个应用程序：", [app["name"] for app in app_list])
    if app_selected:
        if st.button(f"运行 {app_selected}"):
            result = run_app(app_selected)
            if result == 200:
                st.success(f"{app_selected} 已成功运行！")
            else:
                st.error(f"运行 {app_selected} 时出错！返回值 {result}")

# 创建应用程序视图
elif selected_view == "创建应用程序":
    st.subheader("创建新应用程序")
    app_name_input = st.text_input("应用程序名称")
    doip_input = st.text_input("DOIP")
    run_script_input = st.text_area("运行脚本")
    create_app_button = st.button("创建")
    if create_app_button:
        create_app(app_name_input, doip_input, run_script_input)
        st.success(f"已成功创建 {app_name_input} 应用程序！")
