const TargetsTable = () => (
  <section className="mt-10">
    <h3 className="text-lg font-semibold mb-4">Últimos Alvos</h3>
    <table className="w-full text-sm border border-zinc-300">
      <thead className="bg-zinc-200">
        <tr>
          <th className="p-2 text-left">Subdomínio</th>
          <th className="p-2 text-left">Status</th>
          <th className="p-2 text-left">Serviço</th>
          <th className="p-2 text-left">Sensível?</th>
        </tr>
      </thead>
      <tbody>
        <tr className="hover:bg-zinc-100">
          <td className="p-2">admin.exemplo.com</td>
          <td className="p-2 text-green-600">Ativo</td>
          <td className="p-2">HTTP 200</td>
          <td className="p-2 font-bold text-red-500">SIM</td>
        </tr>
        <tr className="hover:bg-zinc-100">
          <td className="p-2">dev.exemplo.com</td>
          <td className="p-2 text-gray-500">Inativo</td>
          <td className="p-2">-</td>
          <td className="p-2">Não</td>
        </tr>
      </tbody>
    </table>
  </section>
);

export default TargetsTable;
