import React, { useState } from "react"
import axiosFetching from "../api/AxiosFetch";
import { useMutation } from "@tanstack/react-query";
import { Navigate } from "react-router";


const Auth: React.FC = () => {
    const [login, setLogin] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [error, setError] =useState<string | null>(null);


    const {mutate, isPending} = useMutation({
        mutationFn: async () => {
            const response = await axiosFetching.post("login/endpoint", {login, password}); // ВОЗМОЖНЯ УГРОЗА НО Я ХЗ. ПОПРАВИТЬ ЭНДПОИНТЫ
            return response.data;
        },
        onSuccess: () => {
            setError(null);
            Navigate({to: "/"})// СЮДА ВСТАВИТЬ Navigate на главную страницу
        },
        onError(err: any /*?*/){
            setError(err.response?.data?.massage || "An error massage");
        }
    });

    const handleSubmit = () => {
        if (!login || !password){
            setError("Fill in correct data");
            return
        } 
        
        mutate()

    }

  return  (
    <div className="flex justify-center items-center min-h-screen bg-gray-100">
      <form
        onSubmit={handleSubmit}
        className="bg-white p-8 rounded-lg shadow-lg w-full max-w-md space-y-4"
      >
        <h2 className="text-2xl font-bold text-center text-gray-800">Authentication form </h2>

        <div>
          <label htmlFor="login" className="block text-sm font-medium text-gray-700">
            Login
          </label>
          <input
            id="login"
            value={login}
            onChange={(e) => setLogin(e.target.value)}
            placeholder="Enter your login"
            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            required
          />
        </div>

        <div>
          <label htmlFor="password" className="block text-sm font-medium text-gray-700">
            Password
          </label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Enter your password"
            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            required
          />
        </div>

        {error && <p className="text-red-500 text-sm">{error}</p>}

        <button
          type="submit"
          disabled={isPending}
          className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          {isPending ? 'Logging in...' : 'Login'}
        </button>
      </form>
    </div>
  )

}

export default Auth