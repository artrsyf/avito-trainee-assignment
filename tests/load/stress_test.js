import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export let options = {
    setupTimeout: '300s',
    stages: [
        { duration: '1m', target: 50 },
        { duration: '1m', target: 75 },
        { duration: '1m', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(99)<50'],
        http_req_failed: ['rate<0.10'],
    },
    noConnectionReuse: true,
};

const BASE_URL = 'http://localhost:8080';
const USERS_COUNT = 10000;
const BATCH_SIZE = 150
const MAX_RETRIES = 3;

const MERCH_ITEMS = [
    "t-shirt", "cup", "book", "pen", "powerbank",
    "hoody", "umbrella", "socks", "wallet", "pink-hoody"
];

export function setup() {
    console.log('🔸 Starting setup: registering users...');
    const users = Array.from({ length: USERS_COUNT }, (_, i) => ({
        username: `user${i + 1}`,
        password: `password${i + 1}`,
    }));

    let completed = 0;
    for (let i = 0; i < users.length; i += BATCH_SIZE) {
        const batch = users.slice(i, i + BATCH_SIZE);
        const requests = batch.map(user => ({
            method: 'POST',
            url: `${BASE_URL}/api/auth`,
            body: JSON.stringify({
                username: user.username,
                password: user.password
            }),
            params: {
                headers: { 'Content-Type': 'application/json' },
                tags: { type: 'setup' },
            },
        }));

        const responses = http.batch(requests);
        
        responses.forEach((res, index) => {
            const user = batch[index];
            if (res.status !== 200) {
                console.warn(`⚠️ Failed to register ${user.username}. Retrying...`);
                for (let attempt = 1; attempt <= MAX_RETRIES; attempt++) {
                    const retryRes = http.post(
                        `${BASE_URL}/api/auth`,
                        JSON.stringify(user),
                        { headers: { 'Content-Type': 'application/json' } }
                    );
                    if (retryRes.status === 200) break;
                    if (attempt === MAX_RETRIES) {
                        throw new Error(`❌ Failed to register ${user.username} after ${MAX_RETRIES} attempts`);
                    }
                    sleep(1);
                }
            }
            completed++;
        });

        console.log(`🔄 Progress: ${completed}/${users.length} users registered`);
    }

    console.log('✅ Setup completed: users registered!');
    return { users };
}

export default function (data) {
    const user = data.users[randomIntBetween(0, USERS_COUNT - 1)];
    const token = authenticate(user);

    const operation = randomIntBetween(1, 2);
    switch (operation) {
        case 1:
            const infoRes = http.get(`${BASE_URL}/api/info`, {
                headers: { Authorization: `Bearer ${token}` },
            });
            check(infoRes, { 'Info status is 200': (r) => r.status === 200 });
            break;

        case 2:
            const item = MERCH_ITEMS[randomIntBetween(0, MERCH_ITEMS.length - 1)];
            const buyRes = http.get(`${BASE_URL}/api/buy/${item}`, {
                headers: { Authorization: `Bearer ${token}` },
            });
            check(buyRes, { 'BuyItem status is 200': (r) => r.status === 200 });
            break;
    }

    sleep(2);
}

function authenticate(user) {
    const res = http.post(
        `${BASE_URL}/api/auth`,
        JSON.stringify(user),
        { headers: { 'Content-Type': 'application/json' } }
    );
    check(res, { 'Auth succeeded': (r) => r.status === 200 });
    return res.json('token');
}