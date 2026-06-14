'use client';
import { useEffect, useState } from "react";
import axios from "axios";
import { Project } from "../lib/types";
import { statuslabel, statusClasses } from "../lib/status";
import { formatDate } from "@/lib/format";


export default function Home() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);


  useEffect(() => {
    async function loadProjects() {
      const apiUrl = 
        process.env.NODE_PUBLIC_API_URL || 'http://localhost:8080'

        try {
          const response = await axios.get<Project[]>(`${apiUrl}/projects`);
          if (Array.isArray(response.data)) {
            setProjects(response.data);
          } else {
            setProjects([]);
            setError('Response tidak berformat array.');
          }
        }catch (err) {
          console.error('Axios fetch error', err);
          setError('Gagal memuat data proyek dari server.');
        } finally {
          setLoading(false);
        }
    } loadProjects();
  }, []);

  return (
    <main className="min-h-screen bg-slate-50 text-slate-900 px-4 py-8 md:px-10">
      <div className="mx-auto max-w-6xl">

        <section className="mb-8 rounded-3xl border border-slate-200 bg-white p-6 shadow-lg">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p className="text-sm uppercase tracking-[0.24em] text-slate-500">
                DevAudit Tracker
              </p>
              <h1 className="mt-2 text-3xl font-semibold text-slate-900">
                Dashboard Audit Ticket
              </h1>
              <p className="mt-2 max-w-2xl text-slate-600">
                Menampilkan status setiap ticket dan progres otomatis dari
                folder kerjaan Anda.
              </p>
            </div>
            {/* <div className="rounded-full bg-slate-100 px-4 py-2 text-sm text-slate-700">
              API: {process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}
            </div> */}
          </div>
        </section>

        {loading && (
          <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <p className="text-slate-700">Memuat project...</p>
          </div>
        )}

        {error && (
          <div className="rounded-3xl border border-rose-200 bg-rose-50 p-6 text-rose-800 shadow-sm">
            <p className="font-medium">Error:</p>
            <p>{error}</p>
          </div>
        )}

        {!loading && !error && projects.length === 0 && (
          <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <p className="text-slate-700">
              Belum ada project terdeteksi. Pastikan backend dan workspace
              sudah berjalan.
            </p>
          </div>
        )}

        {!loading && !error && projects.length > 0 && (
          <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
            {projects.map((project) => (
              <article key={project.TicketNumber}  className="overflow-hidden rounded-3xl border border-slate-200 bg-white p-6 shadow-sm transition hover:-translate-y-1 hover:shadow-md">
                
                <div className="mb-4 flex items-start justify-between gap-3">
                  <div>
                    <p className="text-sm font-medium text-slate-500">
                      Ticket
                    </p>
                    <h2 className="mt-1 text-xl font-semibold text-slate-900">
                      {project.TicketNumber}
                    </h2>
                  </div>
                  <span
                    className={`rounded-full px-3 py-1 text-sm font-semibold ${statusClasses[project.Status] ?? 'bg-slate-100 text-slate-800'}`}
                  >
                    {statuslabel[project.Status] ?? project.Status}
                  </span>
                </div>

                 <div className="space-y-3 text-slate-600">
                  <div className ="rounded-2xl bg-slate-100 p-3">
                    <p className= "text-xs uppercase tracking-[0.18em] text-slate-500">
                      Progress
                    </p>
                    <div className="mt-2 flex gap-2 text-sm">
                      <span>FSD</span>
                      <span className={project.HasFSD ? 'font-semibold text-slate-900' : 'text-slate-400'}>
                        {project.HasFSD ? '✔️' : '✖️'}
                      </span>
                    </div>

                    <div className="mt-2 flex gap-2 text-sm">
                      <span>Analysis</span>
                      <span className={project.HasAnalysis ? 'font-semibold text-slate-900' : 'text-slate-400'}>
                        {project.HasAnalysis ? '✔️' : '✖️'}
                      </span>
                    </div>

                    <div className="mt-1 flex gap-2 text-sm">
                      <span>SIT</span>
                      <span className={project.HasSIT ? 'font-semibold text-slate-900' : 'text-slate-400'}>
                        {project.HasSIT ? '✔️' : '✖️'}
                      </span>
                    </div>

                    <div className="mt-1 flex gap-2 text-sm">
                      <span>Doc Done</span>
                      <span className="font-semibold text-slate-900">
                        {project.DocCount}/4
                      </span>
                    </div>
                  </div>

                  <div className="grid gap-2 text-sm text-slate-500">
                    <div className="flex justify-between">
                      <span>Dibuat</span>
                      <span>{formatDate(project.CreatedAt)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Terakhir update</span>
                      <span>{formatDate(project.UpdatedAt)}</span>
                    </div>
                  </div>
                 </div>


                  <div className="mt-6 flex items-center justify-between">
                    <a href={`/comparison/${encodeURIComponent(project.TicketNumber)}`} className="rounded-2xl bg-slate-900 px-4 py-2 text-sm font-semibold text-white transition hover:bg-slate-800">
                    Lihat Comparison
                    </a>
                  </div>
              </article>
            ))}
          </div>
        )}
      </div>
    </main>
  );
}