import requests

baseURL = 'http://chatbot-backend.mlops.pub/'

chatHeaders = {
    'Content-Type': 'application/json'
}

if __name__ == "__main__":
    authURL = baseURL + 'auth'
    
    authPayload = {
        'username': 'temporary_user',
        'password': 'default_password'
    }

    try:
        response = requests.post(url=authURL, data=authPayload, )
    except Exception as e:
        print(e)
    else:
        # status_code = response.status_code
        access_token = response.json()['accessToken']
        print(access_token)

    chatURL = baseURL + 'chatCompletion'
    chatHeaders['Authorization'] = 'Bearer' + ' ' + access_token
    chatPayload = [
        { "role": "user", "content": "请给我推荐一部喜剧。"},
        { "role": "assistant", "content": "如果你想看一部轻松愉快的喜剧，我推荐你观看《摔跤吧！爸爸》（Dangal）。"},
        { "role": "user", "content": "详细介绍下。"},
        ]

    try:
        response = requests.post(url=chatURL, headers=chatHeaders, json=chatPayload)
    except Exception as e:
        print(e)
    else:
        # print(response)
        print(response.json())