import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Header from "../components/Header";

interface Result {
  domain: string;
  timestamp: string;
  subdomains: number;
  activeSites: number;
  juicyTargets: number;
  path: string;
}

const ResultsPage = () => {
  const [results, setResults] = useState<Result[]>([]);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    fetch("http://localhost:8080/api/results")
      .then((res) => res.json())
      .then((data) => setResults(data))
      .catch(() => setError("Erro ao buscar resultados."));
  }, []);

  return (
    <div className="flex-1">
      <Header />
      <section className="bg-white shadow rounded-lg p-6">
        <h2 className="text-xl font-bold mb-4">Resultados Anteriores</h2>
        {error && <p className="text-red-600">{error}</p>}
        {results.length === 0 && !error && <p>Nenhum resultado encontrado.</p>}
        <ul className="divide-y divide-zinc-200">
          {results.map((r, i) => (
            <li
              key={i}
              onClick={() => navigate(`/resultados/${r.path}`)}
              className="p-4 cursor-pointer hover:bg-zinc-100 rounded-md"
            >
              <div className="font-semibold">{r.domain}</div>
              <div className="text-sm text-zinc-600">
                Data: {r.timestamp} | Subdomínios: {r.subdomains} | Ativos:{" "}
                {r.activeSites} | Juicy: {r.juicyTargets}
              </div>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
};

export default ResultsPage;
