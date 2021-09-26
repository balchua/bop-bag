import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
    vus: 100,
    duration: '30s',
};


export default function () {
    let url = "http://localhost:8000/api/v1/task"

    var payload = JSON.stringify(
    { 
        title: "My First Task", 
        details: "Here you go, this is what i should do", 
        createdDate: "2021-10-25"
    });

    var params = {
        headers: {
          'Content-Type': 'application/json',
        },
      };

    http.post(url, payload, params);
    sleep(1);
}
