## Get user token by request  
1. Get the authrizaiton token by making a request to the auth endpoint, you have two choice here:  

    *  Request with `form-data` :
        ```bash  
        POST /auth HTTP/1.1
        Host: chatbot-backend.mlops.pub
        Content-Type: multipart/form-data

        {
            "username": "temporary_user",
            "password": "default_password"
        }
        ```
    * Requset with `params`:
        ```bash  
        POST /auth HTTP/1.1
        Host: chatbot-backend.mlops.pub
        params: username=temporary_user&password=default_password
        ```

2. The both responses should look similar to:
    ```bash
    HTTP/1.1 200 OK
    Content-Type: application/json

    {
        "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6MSwiaWF0IjoxNDQ0OTE3NjQwLCJuYmYiOjE0NDQ5MTc2NDAsImV4cCI6MTQ0NDkxNzk0MH0.KPmI6WSjRjlpzecPvs3q_T3cJQvAgJvaQAPtk1abC_E"
    }
    ```

3. You may receive errors caused by username/password errors and other issues, the main categories are as follows:
    * The format of the username/password in the request is incorrect or not filled:
        ```bash
        Status:  401 Unauthorized

        {
        "errorMessage": "Authentication failed: cannot find username/password in request body."
        }
        ```
    * The user in the request is not registered in the database or the password is wrong:
        ```bash
        Status:  401 Unauthorized

         {
         "errorMessage": "Authentication failed: cannot authenticate with provided username and password."
         }
        ```

## Request to protected endpoints with token
1. Token received can then be used to make requests against protected `/chatCompletion` or `/chatCompletionStream` endpoints:
    ```bash
    POST /chatCompletion HTTP/1.1
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6MSwiaWF0IjoxNDQ0OTE3NjQwLCJuYmYiOjE0NDQ5MTc2NDAsImV4cCI6MTQ0NDkxNzk0MH0.KPmI6WSjRjlpzecPvs3q_T
    Content-Type: application/json

    [
        {
        "role": "user",
        "Content": "What is your name."
        }
    ]
    ```
    **Notice:** If you need to make requests against protected `/chatCompletionStream` endpoints, make "Accept" of http header be "text/event-stream".

2. Requests to the `chatCompletion` endpoint will get the complete content reponese:
    ```bash
    {
        "id": "chatcmpl-700IlaB0PJwlSijA6hTIAqOXaiE9E",
        "object": "chat.completion",
        "created": 1680235515,
        "model": "gpt-3.5-turbo-0301",
        "choices": [
            {
                "index": 0,
                "message": {
                    "role": "assistant",
                    "content": "Hello, I am XiaoZhi, the openEuler community assistant. How may I assist you today?"
                },
                "finish_reason": "stop"
            }
        ],
        "usage": {
            "prompt_tokens": 24,
            "completion_tokens": 22,
            "total_tokens": 46
        }
    }
    ```
    `choices[0]["message"]["content"]` in the body is the content that needs to be returned by the dialogue.

3. For the request of `chatCompletionStream` endpoint, the relevant content will be returned through Server-Sent Events (SSE). You need to process the EventSource to obtain all the content. The format of each transmission content is:
    ```bash
    event: message
    data:  community
    ```
    `data` item is the content returned each time, which can be directly used to combine for complete content, notice that the data in the first and last event is empty.

    For example, the answer `Hello, XiaoZhi.` from server returns through following results sequentially in stream mode:
    ```bash
    event: message
    data: 

    event: message
    data: Hello

    event: message
    data: ,

    event: message
    data:  Xiao

    event: message
    data: Z

    event: message
    data: hi.

    event: message
    data: 
    ```

4. You may also encounter some error returns:

    * Incorrect token format or not set:
        ```bash
         Status:  401 Unauthorized

         {
         "errorMessage": "Invalid token: empty or not starts with 'Bearer '"
         }
        ```
    * Forged token or has been expired:
        ```bash
         Status:  401 Unauthorized

         {
         "errorMessage": "Invalid token: expired or fake token."
         }
        ```
    * Issues caused by internal programs or network problems on the server side:
        ```bash
         Status:  500 InternalServerError

         {
         "errorMessage": "Failed with chatbot.ChatCompletion: {error}"
         }
        ```

**Notice:** For other network issues, refer to the general Http protocol.  

<br />

## Handle SSE stream in ***vue***
1. Install plugin needed
    ```bash
    npm install @microsoft/fetch-event-source
    ```
2. fetch EventSource
    ```js
    import { fetchEventSource } from '@microsoft/fetch-event-source';

    export function getChatRes (inputText, params) {
    const { messgae } = params
    const headers = {
        'Authorization': 'Bearer' + ' ' + localStorage.getItem('Access-Token') + 5,
    };
    const body = JSON.stringify([
        {
        role: 'user',
        Content: inputText
        }
    ]);
    const es = new fetchEventSource('/chatCompletionStream', {
        method: 'POST',
        headers,
        body,
        async onopen (response) {
        if (response.ok) {
            return; // everything's good
        } else if (response.status >= 400 && response.status < 500 && response.status !== 429) {
            console.log(response.statusText); // handling error
            // throw new Error(response.statusText);
        } else {
            console.log(response.statusText); // handling error
            throw new Error();
        }
        },
        onmessage (event) {
        messgae(event.data);
        },
        onclose () {
        // if the server closes the connection unexpectedly, retry:
        },
        onerror (err) {
        console.log(err)
        throw new Error();
        }
    });
    }
    ```

3. Call the function
    ```js
    getChatRes(this.question, {
            messgae: (res) => {
            console.log(res)
            },
        })
    ```