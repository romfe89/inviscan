import Sidebar from "./components/Sidebar";
import { Routes, Route } from "react-router-dom";
import DashboardPage from "./pages/DashboardPage";
import TargetsPage from "./pages/TargetsPage";

function App() {
  return (
    <div className="flex h-screen bg-zinc-100 text-zinc-800">
      <Sidebar />
      <main className="flex-1 p-8 overflow-y-auto">
        <Routes>
          <Route path="/" element={<DashboardPage />} />
          <Route path="/alvos" element={<TargetsPage />} />
        </Routes>
      </main>
    </div>
  );
}

export default App;
