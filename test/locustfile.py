from locust import HttpUser, task, between
import random

class MyUser(HttpUser):
    wait_time = between(1, 3)

    @task
    def test_endpoint(self):
        mode = random.choice(['restock', 'buy'])
        uid = random.randint(0, 9)
        pid = random.randint(0, 9)
        qty = random.randint(1, 3)

        payload = {
            'mode': mode,
            'uid': uid,
            'pid': pid,
            'qty': qty
        }

        self.client.post("/", json=payload)