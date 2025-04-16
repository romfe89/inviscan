import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import Header from "../components/Header";

interface ResultData {
  subdomains: string[];
  active: string[];
  juicy: string[];
}

const ScanResultDetails = () => {
  const { path } = useParams();
  const [data, setData] = useState<ResultData | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!path) return;
    const fetchData = async () => {
      try {
        const res = await fetch(`http://localhost:8080/api/results/${path}`);
        if (!res.ok) throw new Error("Erro ao buscar resultado");
        const json = await res.json();
        setData(json);
      } catch (err) {
        setError("Erro ao carregar dados." + err);
      }
    };
    fetchData();
  }, [path]);

  return (
    <>
      <Header />
      <section className="bg-white shadow rounded-lg p-6">
        <h2 className="text-xl font-bold mb-4">Detalhes da Varredura</h2>
        {error && <p className="text-red-600">{error}</p>}
        {!data && !error && <p>Carregando...</p>}
        {data && (
          <>
            <div className="mb-6">
              <h3 className="text-lg font-semibold mb-2">Subdomínios</h3>
              <ul className="list-disc list-inside text-sm">
                {data.subdomains.map((s, i) => (
                  <li key={i}>{s}</li>
                ))}
              </ul>
            </div>
            <div className="mb-6">
              <h3 className="text-lg font-semibold mb-2">Sites Ativos</h3>
              <ul className="list-disc list-inside text-sm">
                {data.active.map((s, i) => (
                  <li key={i}>{s}</li>
                ))}
              </ul>
            </div>
            <div>
              <h3 className="text-lg font-semibold mb-2">Juicy Targets</h3>
              <ul className="list-disc list-inside text-sm">
                {data.juicy.map((s, i) => (
                  <li key={i}>{s}</li>
                ))}
              </ul>
            </div>
          </>
        )}
      </section>
    </>
  );
};

export default ScanResultDetails;
