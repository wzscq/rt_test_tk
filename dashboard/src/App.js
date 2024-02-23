import {HashRouter,Routes,Route} from "react-router-dom";
import Dashboard from './pages/Dashboard';

function App() {
  return (
    <>
      <HashRouter>
        <Routes>
          <Route path="/" exact={true} element={<Dashboard/>} />
        </Routes>
      </HashRouter>
    </>
  );
}

export default App;
