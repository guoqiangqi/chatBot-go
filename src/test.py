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
    chatPayload = {
        "role": "user",
        "content": "介绍下你自己？"
    },

    try:
        response = requests.post(url=chatURL, headers=chatHeaders, json=chatPayload)
    except Exception as e:
        print(e)
    else:
        # print(response)
        print(response.json())