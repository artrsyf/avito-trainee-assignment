import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

// --- 🔹 КОНФИГУРАЦИЯ НАГРУЗКИ ---
export const options = {
  stages: [
    { duration: '1m', target: 100 },   // Разогрев до 100 RPS за 1 минуту
    { duration: '3m', target: 1000 },  // Поддержание 1000 RPS в течение 3 минут
    { duration: '1m', target: 0 },     // Плавное снижение нагрузки
],
  thresholds: {
    http_req_duration: ['p(95)<800'], // 95% запросов должны выполняться < 800 мс
    http_req_failed: ['rate<0.05']    // Ошибки < 5%
  }
};

// --- 🔹 1. СОЗДАНИЕ ПОЛЬЗОВАТЕЛЕЙ ---
export function setup() {
  let users = [];
  const totalUsers = 1000; // Создадим 20 пользователей

  for (let i = 1; i <= totalUsers; i++) {
    const username = `testuser_${i}`;
    const password = 'default_pass';

    // 🔹 Авторизация
    let authRes = http.post(
      `${__ENV.API_URL}/api/auth`,
      JSON.stringify({ username, password }),
      { headers: { 'Content-Type': 'application/json' } }
    );

    if (authRes.status === 200) {
      let token = authRes.json('token');
      users.push({ username, token });
    }
    sleep(0.2); // ⏳ Даем серверу обработать запрос
  }

  return { users };
}

// --- 🔹 2. ОСНОВНОЙ ТЕСТ: ПЕРЕВОДЫ ---
export default function(data) {
  const users = data.users;
  
  if (!users || users.length < 2) {
    console.error('❌ Недостаточно пользователей для теста!');
    return;
  }

  const sender = users[__VU % users.length];  
  let receiverIndex = (Math.floor(Math.random() * users.length));  
  if (receiverIndex === (__VU % users.length)) receiverIndex = (receiverIndex + 1) % users.length;
  const receiver = users[receiverIndex];

  const transferRes = http.post(
    `${__ENV.API_URL}/api/sendCoin`,
    JSON.stringify({ toUser: receiver.username, amount: 1 }), // 💰 Уменьшил сумму
    {
      headers: {
        'Authorization': `Bearer ${sender.token}`,
        'Content-Type': 'application/json'
      }
    }
  );

  check(transferRes, { '✅ Transfer success': (r) => r.status === 200 });

  sleep(randomIntBetween(1, 3));
}
