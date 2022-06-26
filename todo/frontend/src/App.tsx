import { FC } from 'react';
import axios from 'axios';

const App: FC = () => {
  const TODO_ID = 'todo';
  const backendUrl = import.meta.env.VITE_BACKEND_URL;

  const submitTodo = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const form = e.currentTarget;
    const data = new FormData(form);
    axios.post(`${backendUrl}/todos`, data.get(TODO_ID)).then(() => {
      console.log('Todo uploaded successfully');
      form.reset();
    });
  };

  return (
    <main className="app">
      <img src={`${backendUrl}/daily-image`} alt="" className="daily-image" />

      <form onSubmit={submitTodo}>
        <input type="text" name={TODO_ID} id={TODO_ID} maxLength={140} />
        <button type="submit">Create TODO</button>
      </form>

      <ul>
        <li>TODO 1</li>
        <li>TODO 2</li>
      </ul>
    </main>
  );
};

export default App;
