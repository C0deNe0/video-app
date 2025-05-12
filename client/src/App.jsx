import {BrowserRouter, Route,Routes} from "react-router-dom"

import CreateRoom from "./components/CreateRoom"
import Room from "./components/room"

function App() {
  return (
    <>
      <BrowserRouter>
          <Routes>
              <Route path="/"  Component={CreateRoom}/>
              <Route path="/room/:roomID" Component={Room}/>
          </Routes>
          
      </BrowserRouter>
    </>
  )
}

export default App
