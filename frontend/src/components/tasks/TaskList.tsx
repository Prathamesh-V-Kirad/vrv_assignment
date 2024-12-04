import { useEffect, useState } from 'react';
import { Task } from '@/types/task';
import { useTaskStore } from '@/lib/tasks';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { format } from 'date-fns';
import { Pencil, Trash2 } from 'lucide-react';

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

interface TaskListProps {
  onEdit: (task: Task) => void;
}

export function TaskList({ onEdit }: TaskListProps) {
  const [tasks, setTasks] = useState<Task[]>([]);
  const { deleteTask } = useTaskStore();

  useEffect(() => {
    // Fetch tasks from the API
    const fetchTasks = async () => {
      try {
        const response = await fetch('http://localhost:8000/api/tasks', {
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
        });
        if (!response.ok) {
          throw new Error('Failed to fetch tasks');
        }
        const data = await response.json();
        const formattedTasks = data.tasks.map((task: any) => ({
          id: task._id,
          name: task.name,
          description: task.description,
          status: task.status,
          createdAt: new Date(task.created_at),
          updatedAt: new Date(task.updated_at),
        }));

        setTasks(formattedTasks);
      } catch (error) {
        console.error('Error fetching tasks:', error);
      }
    };

    fetchTasks();
  }, []);

  const handleToggleComplete = async (taskId: string, currentStatus: boolean) => {
    try {
      const response = await fetch(`http://localhost:8000/api/tasks/${taskId}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ status: !currentStatus }),
      });

      if (!response.ok) {
        throw new Error('Failed to update task status');
      }

      setTasks((prevTasks) =>
        prevTasks.map((task) =>
          task.id === taskId ? { ...task, status: !currentStatus } : task
        )
      );
    } catch (error) {
      console.error('Error updating task status:', error);
    }
  };

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-12">Status</TableHead>
            <TableHead>Task</TableHead>
            <TableHead>Created</TableHead>
            <TableHead className="w-24">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {tasks.length > 0 ? (
            tasks.map((task) => (
              <TableRow key={task.id}>
                <TableCell>
                  <Checkbox
                    checked={task.status}
                    onCheckedChange={() => handleToggleComplete(task.id, task.status)}
                  />
                </TableCell>
                <TableCell>
                  <div>
                    <p className={task.status ? 'line-through text-gray-500' : ''}>
                      {task.name}
                    </p>
                    <p className="text-sm text-gray-500">{task.description}</p>
                  </div>
                </TableCell>
                <TableCell className="text-sm text-gray-500">
                  {format(task.createdAt ?? new Date(), 'PPp')}
                </TableCell>
                <TableCell>
                  <div className="flex space-x-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => onEdit(task)}
                      className="h-8 w-8"
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => deleteTask(task.id)}
                      className="h-8 w-8 text-destructive"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={4} className="text-center py-8 text-gray-500">
                No tasks yet. Create one to get started!
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
