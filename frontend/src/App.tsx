import Sidebar from "./components/Sidebar";
import { Routes, Route } from "react-router-dom";
import DashboardPage from "./pages/DashboardPage";
import TargetsPage from "./pages/TargetsPage";
import ScanPage from "./pages/ScanPage";
import ResultsPage from "./pages/ResultsPage";
import ScanResultDetails from "./pages/ScanResultDetails";

function App() {
  return (
    <div className="flex h-screen bg-zinc-100 text-zinc-800">
      <Sidebar />
      <main className="flex-1 p-8 overflow-y-auto">
        <Routes>
          <Route path="/" element={<DashboardPage />} />
          <Route path="/alvos" element={<TargetsPage />} />
          <Route path="/varreduras" element={<ScanPage />} />
          <Route path="/resultados" element={<ResultsPage />} />
          <Route path="/resultados/:path" element={<ScanResultDetails />} />
        </Routes>
      </main>
    </div>
  );
}

export default App;
