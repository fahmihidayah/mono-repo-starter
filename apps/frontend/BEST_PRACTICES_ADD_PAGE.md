# Best Practices for Add/Create Pages

A comprehensive guide for building form pages in React Router v7 with TypeScript, Zod, and shadcn/ui.

## Table of Contents
1. [File Structure & Routing](#file-structure--routing)
2. [Form Implementation](#form-implementation)
3. [Validation Strategy](#validation-strategy)
4. [Error Handling](#error-handling)
5. [User Experience](#user-experience)
6. [Security & Performance](#security--performance)

---

## File Structure & Routing

### ✅ Recommended File Naming

```
app/routes/
├── dashboard.tsx                    # Layout with <Outlet />
├── dashboard.tasks.tsx              # List page
└── dashboard.tasks.add.tsx          # Add page (nested)
```

**Why this structure?**
- Creates URL: `/dashboard/tasks/add`
- Inherits dashboard layout (sidebar, header)
- Logical grouping by feature (tasks)

### ❌ Avoid

```
app/routes/
└── dashboard.add.tsx               # Too generic, unclear what you're adding
```

---

## Form Implementation

### ✅ Best Practice: Server Actions with Client Validation

```tsx
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form as ReactRouterForm } from "react-router";

// 1. Define schema (shared between client and server)
const taskSchema = z.object({
  title: z.string().min(1).max(100),
  description: z.string().min(10).max(500),
});

type TaskFormData = z.infer<typeof taskSchema>;

// 2. Server action with validation
export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();

  const data = {
    title: formData.get('title')?.toString() || '',
    description: formData.get('description')?.toString() || '',
  };

  // Server-side validation (CRITICAL!)
  try {
    const validated = taskSchema.parse(data);

    // Insert to database
    await db.insert(task).values({
      id: randomUUID(),
      ...validated,
      createdAt: new Date(),
    });

    return redirect("/dashboard/tasks");
  } catch (error) {
    if (error instanceof z.ZodError) {
      return {
        success: false,
        error: "Validation failed",
        fieldErrors: error.flatten().fieldErrors
      };
    }

    return { success: false, error: "Server error" };
  }
}

// 3. Component with React Hook Form
export default function AddPage() {
  const form = useForm<TaskFormData>({
    resolver: zodResolver(taskSchema),
    defaultValues: { title: "", description: "" }
  });

  return (
    <Form {...form}>
      <ReactRouterForm method="post">
        {/* Form fields */}
      </ReactRouterForm>
    </Form>
  );
}
```

### Key Points

✅ **ALWAYS validate on the server** - Never trust client input
✅ **Use the same schema** - Client and server use identical validation
✅ **Type safety** - Infer types from Zod schema
✅ **Progressive enhancement** - Works without JavaScript

---

## Validation Strategy

### Server-Side (REQUIRED)

```tsx
export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();

  // Extract data
  const data = Object.fromEntries(formData);

  // Validate with Zod
  const result = taskSchema.safeParse(data);

  if (!result.success) {
    return {
      success: false,
      errors: result.error.flatten().fieldErrors
    };
  }

  // Safe to use result.data
  await db.insert(task).values(result.data);
  return redirect("/success");
}
```

### Client-Side (OPTIONAL but recommended)

```tsx
const form = useForm<TaskFormData>({
  resolver: zodResolver(taskSchema), // Validates before submit
  defaultValues: { title: "", description: "" }
});
```

**Benefits:**
- Immediate feedback to users
- Reduces server load
- Better UX

---

## Error Handling

### ✅ Comprehensive Error Handling

```tsx
export async function action({ request }: Route.ActionArgs) {
  try {
    const formData = await request.formData();
    const data = Object.fromEntries(formData);

    // Validation error
    const validated = taskSchema.parse(data);

    // Database operation
    await db.insert(task).values(validated);

    return redirect("/dashboard/tasks");

  } catch (error) {
    // Validation errors
    if (error instanceof z.ZodError) {
      return {
        success: false,
        error: "Please check your input",
        fieldErrors: error.flatten().fieldErrors
      };
    }

    // Database errors
    if (error.code === '23505') { // Unique constraint
      return {
        success: false,
        error: "A task with this title already exists"
      };
    }

    // Generic errors
    console.error("Task creation failed:", error);
    return {
      success: false,
      error: "Something went wrong. Please try again."
    };
  }
}
```

### Display Errors to User

```tsx
export default function AddPage() {
  const actionData = useActionData<typeof action>();

  // Show toast on error
  useEffect(() => {
    if (actionData && !actionData.success) {
      toast.error(actionData.error);
    }
  }, [actionData]);

  // Or show inline
  return (
    <div>
      {actionData?.error && (
        <Alert variant="destructive">
          <AlertDescription>{actionData.error}</AlertDescription>
        </Alert>
      )}
      {/* Form */}
    </div>
  );
}
```

---

## User Experience

### 1. Loading States

```tsx
const navigation = useNavigation();
const isSubmitting = navigation.state === "submitting";

<Button type="submit" disabled={isSubmitting}>
  {isSubmitting && <Loader2 className="mr-2 size-4 animate-spin" />}
  {isSubmitting ? "Creating..." : "Create Task"}
</Button>
```

### 2. Success Feedback

**Option A: Toast on redirect (Recommended)**

```tsx
// After successful creation, show toast
import { toast } from "sonner";

// In the list page loader or component
useEffect(() => {
  const params = new URLSearchParams(window.location.search);
  if (params.get('created') === 'true') {
    toast.success("Task created successfully!");
    // Clean URL
    params.delete('created');
    navigate({ search: params.toString() }, { replace: true });
  }
}, []);

// In action, redirect with query param
return redirect("/dashboard/tasks?created=true");
```

**Option B: Flash messages (React Router way)**

Use session storage for flash messages.

### 3. Cancel Button

```tsx
const navigate = useNavigate();

<Button
  type="button"
  variant="outline"
  onClick={() => navigate("/dashboard/tasks")}
>
  Cancel
</Button>
```

### 4. Back Navigation

```tsx
import { ArrowLeft } from "lucide-react";

<Button
  variant="ghost"
  size="sm"
  onClick={() => navigate(-1)}
  className="gap-2"
>
  <ArrowLeft className="size-4" />
  Back
</Button>
```

### 5. Form Layout

```tsx
<div className="flex-1 p-6">
  <div className="max-w-2xl mx-auto space-y-6">
    {/* Back button */}
    <Button variant="ghost" onClick={() => navigate(-1)}>
      <ArrowLeft /> Back
    </Button>

    {/* Form card */}
    <Card>
      <CardHeader>
        <CardTitle>Add New Task</CardTitle>
        <CardDescription>
          Create a new task to organize your work
        </CardDescription>
      </CardHeader>
      <CardContent>
        {/* Form */}
      </CardContent>
    </Card>
  </div>
</div>
```

---

## Security & Performance

### Security Checklist

✅ **Server-side validation** - Always validate on server
✅ **Sanitize input** - Remove dangerous characters
✅ **CSRF protection** - React Router handles this
✅ **SQL injection** - Use parameterized queries (Drizzle does this)
✅ **Rate limiting** - Limit submissions per user
✅ **Authentication** - Verify user is logged in

### Performance Tips

✅ **Client-side validation** - Reduce server requests
✅ **Debounce** - For real-time validation
✅ **Optimistic UI** - Show success immediately
✅ **Progressive enhancement** - Works without JS

---

## Complete Example

See [dashboard.tasks.add.improved.tsx](app/routes/dashboard.tasks.add.improved.tsx) for a full implementation with:

- ✅ Server-side validation
- ✅ Client-side validation
- ✅ Loading states
- ✅ Error handling
- ✅ Success feedback
- ✅ Cancel button
- ✅ Back navigation
- ✅ Proper layout
- ✅ Type safety
- ✅ Accessibility

---

## Comparison: Current vs. Best Practice

### Current Implementation Issues

❌ No server-side validation
❌ Error not displayed to user
❌ No loading state
❌ No success feedback
❌ No cancel button
❌ Unused imports

### Best Practice Implementation

✅ Server-side validation with Zod
✅ Error toasts and inline messages
✅ Loading state with spinner
✅ Success toast on redirect
✅ Cancel and back buttons
✅ Clean imports
✅ Better UX with descriptions
✅ Responsive layout

---

## Quick Checklist

When creating an add/create page, ensure:

- [ ] File named correctly (e.g., `feature.add.tsx`)
- [ ] Zod schema defined and exported
- [ ] Server action validates with schema
- [ ] Client uses React Hook Form with Zod resolver
- [ ] Loading state shown during submission
- [ ] Errors displayed to user (toast or inline)
- [ ] Success feedback on completion
- [ ] Cancel button to go back
- [ ] Form disabled during submission
- [ ] Proper redirect after success
- [ ] Type-safe with TypeScript
- [ ] Accessible (labels, ARIA attributes)

---

## Additional Resources

- [React Router v7 Docs](https://reactrouter.com)
- [Zod Documentation](https://zod.dev)
- [React Hook Form](https://react-hook-form.com)
- [shadcn/ui Components](https://ui.shadcn.com)
