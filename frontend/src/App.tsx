import React, { useState } from "react";
import axios from "axios";
import './App.css';

const App: React.FC = () => {
  const [number, setNumber] = useState<string>("");
  const [response, setResponse] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async () => {
    if (!number) {
      setError("Please enter a number.");
      setResponse(null);
      return;
    }

    try {
      // Clear previous errors
      setError(null);

      // Send request to the backend
      const res = await axios.get(`http://localhost:8080/search/${number}`);
      setResponse(res.data); // Expecting JSON from the backend
    } catch (err: any) {
      if (err.response && err.response.status === 404) {
        setError("Number not found.");
      } else if (err.response && err.response.status === 400) {
        setError("Invalid input. Please enter a valid number.");
      } else {
        setError("An error occurred. Please try again.");
      }
      setResponse(null);
    }
  };

  return (
    <div className="App">
      <h1 className="Search-header" >Search Number</h1>
      <input
        type="number"
        className="Search-input"
        value={number}
        onChange={(e) => setNumber(e.target.value)}
        placeholder="Enter a number"
      />
      <button
        onClick={handleSearch}
        className="Search-button"
      >
        Search
      </button>
      {response && (
        <div className="Search-success">
          <strong>Index:</strong> {JSON.stringify(response)}
        </div>
      )}
      {error && (
        <div className="Search-error">
          {error}
        </div>
      )}
    </div>
  );
};

export default App;
