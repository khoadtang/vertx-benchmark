import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
  stages: [
    { duration: '1m', target: 1000 } // Ramp up to 1000 RPS over 1 minute
  ],
};

export default function () {
  // Send a GET request to the specified endpoint
  const response = http.get('http://localhost:8080/vertx/fetch');

  // Print the response status to the console
  console.log(`Response status: ${response.status}`);

  // Sleep for a random duration between 0.5s and 1s
  sleep(Math.random() * (1 - 0.5) + 0.5);
}
