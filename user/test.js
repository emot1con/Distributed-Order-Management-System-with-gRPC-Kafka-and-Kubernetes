import http from "k6/http";
import { check } from "k6";

export const options = {
  scenarios: {
    constant_rps: {
      executor: "constant-arrival-rate",
      rate: 100, // 100 requests per second
      timeUnit: "1s", // Dalam 1 detik
      duration: "5s", // Jalankan selama 30 detik
      preAllocatedVUs: 50, // Minimal 50 virtual users
      maxVUs: 200, // Maksimal 200 virtual users
    },
  },
};

export default function () {
  const url = "http://localhost:8080/auth/login"; // Ganti dengan endpoint auth service kamu
  const payload = JSON.stringify({
    email: "me9@me.com",
    password: "secret",
  });

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  const res = http.post(url, payload, params);

  check(res, {
    "status is 200": (r) => r.status === 200,
    "token received": (r) => r.json("token") !== "",
  });
}
