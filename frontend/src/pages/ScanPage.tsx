import { useState } from "react";
import Header from "../components/Header";

const ScanPage = () => {
  const [url, setUrl] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const handleScan = async () => {
    if (!url) {
      setMessage("Digite uma URL válida.");
      return;
    }

    setLoading(true);
    setMessage("");

    try {
      const response = await fetch("http://localhost:8080/api/scan", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url }),
      });

      if (!response.ok) {
        throw new Error("Falha na requisição");
      }

      const data = await response.json();
      setMessage(data.message);
    } catch (error) {
      setMessage(`Erro ao iniciar a varredura: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Header />
      <section className="bg-white shadow rounded-lg p-6">
        <h3 className="text-lg font-semibold mb-4">Nova Varredura</h3>
        <div className="flex gap-4">
          <input
            type="text"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="Digite a URL do alvo..."
            className="flex-1 px-4 py-2 border border-zinc-300 rounded-md"
          />
          <button
            onClick={handleScan}
            disabled={loading}
            className={`px-4 py-2 rounded-md shadow text-white ${
              loading ? "bg-gray-400" : "bg-blue-600 hover:bg-blue-700"
            }`}
          >
            {loading ? "Escaneando..." : "Escanear"}
          </button>
        </div>
        {message && <p className="mt-4 text-sm text-zinc-600">{message}</p>}
      </section>
    </>
  );
};

export default ScanPage;
