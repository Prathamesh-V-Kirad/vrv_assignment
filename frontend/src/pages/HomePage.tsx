import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { CheckCircle, ListTodo } from 'lucide-react';

export function HomePage() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-gray-100 w-full">
      <div className="mx-5 px-4 py-16">
        <div className="max-w-3xl mx-auto text-center">
          <div className="flex justify-center mb-8">
            <ListTodo className="h-16 w-16 text-primary" />
          </div>
          <h1 className="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl mb-6">
            Manage your tasks with ease
          </h1>
          <p className="text-lg text-gray-600 mb-8">
            Stay organized and boost your productivity with our simple yet powerful
            task management solution.
          </p>
          <div className="flex flex-col sm:flex-row justify-center gap-4">
            <Button asChild size="lg">
              <Link to="/signup">Get Started</Link>
            </Button>
            <Button asChild variant="outline" size="lg">
              <Link to="/login">Login</Link>
            </Button>
          </div>
        </div>

        <div className="mt-20 grid gap-8 md:grid-cols-3">
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <CheckCircle className="h-12 w-12 text-primary mb-4" />
            <h3 className="text-xl font-semibold mb-2">Easy to Use</h3>
            <p className="text-gray-600">
              Intuitive interface that helps you focus on what matters most.
            </p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <CheckCircle className="h-12 w-12 text-primary mb-4" />
            <h3 className="text-xl font-semibold mb-2">Stay Organized</h3>
            <p className="text-gray-600">
              Keep track of all your tasks in one centralized location.
            </p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <CheckCircle className="h-12 w-12 text-primary mb-4" />
            <h3 className="text-xl font-semibold mb-2">Boost Productivity</h3>
            <p className="text-gray-600">
              Complete tasks efficiently and achieve your goals faster.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}