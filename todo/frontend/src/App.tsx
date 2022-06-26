import { FC } from 'react';

const App: FC = () => {
  return (
    <main className="app">
      <img
        src="http://localhost:8081/daily-image"
        alt=""
        className="daily-image"
      />

      <form onSubmit={e => e.preventDefault()}>
        <input type="text" id="todo" maxLength={140} />
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
