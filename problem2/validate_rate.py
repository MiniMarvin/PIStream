import requests
import threading

api_url = 'https://api.pi.delivery/v1/pi'

success = 0
failure = 0
query_size = 1000

def query_pi(start):
    global success
    global failure
    ok = False
    params = {
        'start': start,
        'numberOfDigits': query_size
    }
    try:
        ans = requests.get(api_url, params)
        jsonAns = ans.json()
        if len(jsonAns['content']) == query_size:
            success += 1
            ok = True
    except:
        pass
    
    if not ok:
        failure += 1

if __name__ == "__main__":
    thread_list = []
    for i in range(100000):
        thread = threading.Thread(target=lambda : query_pi(i*1000))
        thread.start()
        thread_list.append(thread)
    
    for t in thread_list:
        t.join()
    
    print("success: ", success, "/", success + failure)