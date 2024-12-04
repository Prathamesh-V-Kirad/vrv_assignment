import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { TaskList } from '@/components/tasks/TaskList';
import { TaskForm } from '@/components/tasks/TaskForm';
import { Task } from '@/types/task';
import { Header } from '@/components/layout/Header';
import { Plus } from 'lucide-react';

export function TasksPage() {
  const [open, setOpen] = useState(false);
  const [selectedTask, setSelectedTask] = useState<Task | undefined>();

  const handleEdit = (task: Task) => {
    setSelectedTask(task);
    setOpen(true);
  };

  const handleOpenChange = (open: boolean) => {
    setOpen(open);
    if (!open) {
      setSelectedTask(undefined);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <main className="mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-8">
          <h2 className="text-3xl font-bold text-gray-900">Tasks</h2>
          <Button onClick={() => setOpen(true)}>
            <Plus className="mr-2 h-4 w-4" /> Add Task
          </Button>
        </div>
        <TaskList onEdit={handleEdit} />
        <TaskForm
          task={selectedTask}
          open={open}
          onOpenChange={handleOpenChange}
        />
      </main>
    </div>
  );
}