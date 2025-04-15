import { useEffect, useState } from "react";

interface ScanResult {
  domain: string;
  timestamp: string;
  subdomains: number;
  activeSites: number;
  juicyTargets: number;
  path: string;
}

const ResultsPage = () => {
  const [results, setResults] = useState<ScanResult[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("http://localhost:8080/api/results")
      .then((res) => res.json())
      .then((data) => setResults(data))
      .catch(() => setResults([]))
      .finally(() => setLoading(false));
  }, []);

  return (
    <section className="bg-white shadow rounded-lg p-6">
      <h3 className="text-lg font-semibold mb-4">Histórico de Varreduras</h3>
      {loading ? (
        <p className="text-sm text-zinc-500">Carregando...</p>
      ) : (
        <table className="w-full text-sm border border-zinc-300">
          <thead className="bg-zinc-200">
            <tr>
              <th className="p-2 text-left">Domínio</th>
              <th className="p-2 text-left">Data</th>
              <th className="p-2 text-left">Subdomínios</th>
              <th className="p-2 text-left">Ativos</th>
              <th className="p-2 text-left">Juicy</th>
            </tr>
          </thead>
          <tbody>
            {results.map((res, idx) => (
              <tr key={idx} className="hover:bg-zinc-100">
                <td className="p-2">{res.domain}</td>
                <td className="p-2">{res.timestamp}</td>
                <td className="p-2">{res.subdomains}</td>
                <td className="p-2">{res.activeSites}</td>
                <td className="p-2">{res.juicyTargets}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  );
};

export default ResultsPage;
