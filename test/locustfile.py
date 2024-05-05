from locust import HttpUser, task, between, LoadTestShape
import random
import math

class MyUser(HttpUser):
    wait_time = between(1, 3)

    @task
    def test_endpoint(self):
        mode = random.choice(['restock', 'buy'])
        uid = random.randint(0, 9)
        pid = random.randint(0, 9)
        qty = random.randint(1, 3)

        params = {
            'mode': mode,
            'uid': uid,
            'pid': pid,
            'qty': qty
        }

        self.client.post("/", params=params)

class SineLoadShape(LoadTestShape):
    period = 60
    min_users = 10
    max_users = 50

    def tick(self):
        run_time = self.get_run_time()
        user_range = self.max_users - self.min_users
        amplitude = user_range / 2
        average_users = self.min_users + amplitude
        user_count = round(average_users + amplitude * math.sin(2 * math.pi * run_time / self.period))
        spawn_rate = user_count - self.get_current_user_count()
        if spawn_rate == 0:
            spawn_rate = 1
        return user_count, spawn_rate
