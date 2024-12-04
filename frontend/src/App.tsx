import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Toaster } from '@/components/ui/toaster';
import { HomePage } from '@/pages/HomePage';
import { SignUpPage } from '@/pages/SignUpPage';
import { LoginPage } from '@/pages/LoginPage';
import { TasksPage } from '@/pages/TasksPage';
import { PrivateRoute } from '@/components/PrivateRoute';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/signup" element={<SignUpPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route
          path="/tasks"
          element={
            <PrivateRoute>
              <TasksPage />
            </PrivateRoute>
          }
        />
      </Routes>
      <Toaster />
    </Router>
  );
}

export default App;
