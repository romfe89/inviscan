import { NavLink } from "react-router-dom";

const Sidebar = () => (
  <aside className="w-60 bg-zinc-900 text-white p-6 flex flex-col justify-between">
    <div>
      <h1 className="text-2xl font-bold mb-8">Inviscan</h1>
      <nav className="flex flex-col gap-4">
        <NavLink
          to="/"
          className={({ isActive }) =>
            isActive ? "text-blue-400" : "hover:text-blue-400"
          }
        >
          Dashboard
        </NavLink>
        <NavLink
          to="/alvos"
          className={({ isActive }) =>
            isActive ? "text-blue-400" : "hover:text-blue-400"
          }
        >
          Alvos
        </NavLink>
        <NavLink
          to="/varreduras"
          className={({ isActive }) =>
            isActive ? "text-blue-400" : "hover:text-blue-400"
          }
        >
          Varreduras
        </NavLink>
        <NavLink
          to="/resultados"
          className={({ isActive }) =>
            isActive ? "text-blue-400" : "hover:text-blue-400"
          }
        >
          Resultados
        </NavLink>
        <a href="#">Sobre</a>
      </nav>
    </div>
    <footer className="text-xs text-gray-400 text-center mt-8">
      Inviscan © 2025
    </footer>
  </aside>
);

export default Sidebar;
