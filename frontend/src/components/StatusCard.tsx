import { useState, useEffect } from "react";

const StatusCard = () => {
  const [ping, setPing] = useState("...");

  useEffect(() => {
    fetch("http://localhost:8080/api/ping")
      .then((res) => res.json())
      .then((data) => setPing(data.message))
      .catch(() => setPing("Offline"));
  }, []);

  return (
    <section className="bg-white shadow rounded-lg p-6">
      <h3 className="text-lg font-semibold mb-4">Status do Backend</h3>
      <p>
        Resposta da API: <strong>{ping}</strong>
      </p>
    </section>
  );
};

export default StatusCard;
