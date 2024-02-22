import {HashRouter,Routes,Route} from "react-router-dom";
import Monitor from './pages/Monitor';

function App() {
  return (
    <>
      <HashRouter>
        <Routes>
          <Route path="/" exact={true} element={<Monitor/>} />
        </Routes>
      </HashRouter>
    </>
  );
}

export default App;
