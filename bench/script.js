import http from 'k6/http';
import { check, group, sleep } from 'k6';

export let options = {
  vus: 100,
  duration: '30s',
};


export default function () {
  let add_url = "http://localhost:8000/api/v1/task"
  let get_url = "http://localhost:8000/api/v1/task"
  let get_all_url = "http://localhost:8081/api/v1/tasks"

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

  group('simple add task and query', (_) => {
    // Add task
    let add_response = http.post(add_url, payload, params);
    check(add_response, {
      'is status 200': (r) => r.status === 200,
      'is id present': (r) => r.json().hasOwnProperty('id'),
    });
    var id = add_response.json()['id'];
    // Get task by id request
    let get_task_response = http.get(
      get_url + "/" + id,
      params,
    );

  });
}

