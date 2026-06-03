import { useState, useEffect } from "react";
import "./App.css";

function App() {
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // URL backend (thay bằng URL thật của bạn)
  const API_URL = "https://recruit-be.herokuapp.com/api/jobs";

  useEffect(() => {
    fetch(API_URL)
      .then((res) => {
        if (!res.ok) throw new Error("Network error");
        return res.json();
      })
      .then((data) => {
        setJobs(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, []);

  return (
    <div className="App">
      <header className="hero">
        <h1>🚀 Cơ Hội Việc Làm IT AIAIAIAI AIAIAIAI</h1>
        <p>Khám phá công việc mơ ước tại các công ty hàng đầu</p>
      </header>

      <section className="jobs-section">
        <h2>Vị trí đang tuyển</h2>

        {loading && <p className="loading">Đang tải...</p>}
        {error && <p className="error">Lỗi: {error}</p>}

        <div className="jobs-grid">
          {jobs.map((job) => (
            <div key={job.id} className="job-card">
              <h3>{job.title}</h3>
              <p className="location">📍 {job.location}</p>
              <p className="description">{job.description}</p>
              <button className="apply-btn">Ứng tuyển ngay</button>
            </div>
          ))}
          {!loading && jobs.length === 0 && (
            <p>Chưa có công việc nào được đăng.</p>
          )}
        </div>
      </section>
    </div>
  );
}

export default App;