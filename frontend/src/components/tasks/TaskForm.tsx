import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useTaskStore } from '@/lib/tasks';
import { Task } from '@/types/task';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';

const taskSchema = z.object({
  name: z.string().min(1, 'Title is required'),
  description: z.string().min(1, 'Description is required'),
  status: z.boolean(),
});

interface TaskFormProps {
  task?: Task;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function TaskForm({ task, open, onOpenChange }: TaskFormProps) {
  const { addTask, updateTask } = useTaskStore();

  const form = useForm<z.infer<typeof taskSchema>>({
    resolver: zodResolver(taskSchema),
    defaultValues: {
      name: task?.name || '',  // Pre-fill with existing task data
      description: task?.description || '',
      status: task?.status || false,
    },
  });

  const onSubmit = async (values: z.infer<typeof taskSchema>) => {
    if (task) {
      await handleUpdateTask(task.id, {
        name: values.name,
        description: values.description,
        status: values.status,
      });
    } else {
      await handleAddTask({
        name: values.name,
        description: values.description,
        status: values.status,
      });
    }
    onOpenChange(false);  // Close the form after submission
    form.reset();  // Reset the form after submission
  };

  const handleAddTask = async ({
    name,
    description,
    status,
  }: { name: string; description: string; status: boolean }) => {
    try {
      const response = await fetch('http://localhost:8000/api/tasks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ name, description, status }),
      });

      if (!response.ok) {
        throw new Error('Failed to add task');
      }

      const data = await response.json();
      console.log('Task added:', data);

      // Add the new task to the task list (without needing to re-fetch)
      addTask(data.task);
    } catch (error) {
      console.error('Error adding task:', error);
    }
  };

  const handleUpdateTask = async (
    taskId: string,
    { name, description, status }: { name: string; description: string; status: boolean }
  ) => {
    try {
      const response = await fetch(`http://localhost:8000/api/tasks/${taskId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ name, description, status }),
      });

      if (!response.ok) {
        throw new Error('Failed to update task');
      }

      const data = await response.json();
      console.log('Task updated:', data);
      updateTask(taskId, data.task); // Update the task list state
    } catch (error) {
      console.error('Error updating task:', error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{task ? 'Edit Task' : 'Create Task'}</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Title</FormLabel>
                  <FormControl>
                    <Input placeholder="Task title" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Task description"
                      className="resize-none"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="flex justify-end space-x-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
              >
                Cancel
              </Button>
              <Button type="submit">{task ? 'Update' : 'Create'}</Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
