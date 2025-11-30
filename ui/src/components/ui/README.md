# UI Components (Shadcn/ui)

## Overview

This folder contains primitive UI components from [Shadcn/ui](https://ui.shadcn.com/), built on [Radix UI](https://www.radix-ui.com/) primitives with Tailwind CSS styling. These components are accessible, customizable, and serve as building blocks for feature components.

## Why Shadcn/ui?

- **Accessible**: Built on Radix UI with full keyboard and screen reader support
- **Customizable**: Copy-paste components you can modify directly
- **Unstyled**: Use your own design system with Tailwind
- **No runtime**: Components are part of your codebase, not a dependency

## Components (26)

### Form Controls

| Component | File | Purpose | Radix UI |
|-----------|------|---------|----------|
| Button | `button.tsx` | Clickable actions | Slot |
| Input | `input.tsx` | Text input | - |
| Textarea | `textarea.tsx` | Multi-line input | - |
| Select | `select.tsx` | Dropdown selection | Select |
| Checkbox | `checkbox.tsx` | Binary selection | Checkbox |
| Radio Group | `radio-group.tsx` | Single selection | RadioGroup |
| Switch | `switch.tsx` | Toggle switch | Switch |
| Label | `label.tsx` | Form labels | Label |

### Display

| Component | File | Purpose | Radix UI |
|-----------|------|---------|----------|
| Card | `card.tsx` | Content container | - |
| Badge | `badge.tsx` | Status indicators | - |
| Avatar | `avatar.tsx` | User images | Avatar |
| Alert | `alert.tsx` | Information banners | - |
| Progress | `progress.tsx` | Progress bars | Progress |
| Skeleton | `skeleton.tsx` | Loading placeholders | - |
| Separator | `separator.tsx` | Visual dividers | Separator |
| Table | `table.tsx` | Data tables | - |

### Overlay

| Component | File | Purpose | Radix UI |
|-----------|------|---------|----------|
| Dialog | `dialog.tsx` | Modal dialogs | Dialog |
| Alert Dialog | `alert-dialog.tsx` | Confirmation dialogs | AlertDialog |
| Sheet | `sheet.tsx` | Slide-out panels | Dialog |
| Popover | `popover.tsx` | Floating content | Popover |
| Dropdown Menu | `dropdown-menu.tsx` | Action menus | DropdownMenu |

### Navigation

| Component | File | Purpose | Radix UI |
|-----------|------|---------|----------|
| Tabs | `tabs.tsx` | Tab navigation | Tabs |
| Pagination | `pagination.tsx` | Page navigation | - |
| Scroll Area | `scroll-area.tsx` | Custom scrollbars | ScrollArea |

### Feedback

| Component | File | Purpose | Radix UI |
|-----------|------|---------|----------|
| Toast | `toast.tsx` | Notifications | Toast |
| Toaster | `toaster.tsx` | Toast container | - |

## Usage Examples

### Button

```tsx
import { Button } from "@/components/ui/button";

// Variants
<Button variant="default">Primary</Button>
<Button variant="secondary">Secondary</Button>
<Button variant="destructive">Danger</Button>
<Button variant="outline">Outline</Button>
<Button variant="ghost">Ghost</Button>
<Button variant="link">Link</Button>

// Sizes
<Button size="sm">Small</Button>
<Button size="default">Default</Button>
<Button size="lg">Large</Button>
<Button size="icon"><Icon /></Button>

// States
<Button disabled>Disabled</Button>
<Button asChild><Link href="/page">As Link</Link></Button>
```

### Card

```tsx
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";

<Card>
  <CardHeader>
    <CardTitle>Card Title</CardTitle>
    <CardDescription>Card description text</CardDescription>
  </CardHeader>
  <CardContent>
    <p>Card content goes here.</p>
  </CardContent>
  <CardFooter>
    <Button>Action</Button>
  </CardFooter>
</Card>
```

### Dialog

```tsx
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";

<Dialog>
  <DialogTrigger asChild>
    <Button>Open Dialog</Button>
  </DialogTrigger>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Dialog Title</DialogTitle>
      <DialogDescription>
        This is a description of the dialog content.
      </DialogDescription>
    </DialogHeader>
    <div>Dialog body content</div>
    <DialogFooter>
      <Button variant="outline">Cancel</Button>
      <Button>Confirm</Button>
    </DialogFooter>
  </DialogContent>
</Dialog>
```

### Select

```tsx
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";

<Select value={value} onValueChange={setValue}>
  <SelectTrigger>
    <SelectValue placeholder="Select option" />
  </SelectTrigger>
  <SelectContent>
    <SelectItem value="option1">Option 1</SelectItem>
    <SelectItem value="option2">Option 2</SelectItem>
    <SelectItem value="option3">Option 3</SelectItem>
  </SelectContent>
</Select>
```

### Tabs

```tsx
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";

<Tabs defaultValue="tab1">
  <TabsList>
    <TabsTrigger value="tab1">Tab 1</TabsTrigger>
    <TabsTrigger value="tab2">Tab 2</TabsTrigger>
  </TabsList>
  <TabsContent value="tab1">Tab 1 content</TabsContent>
  <TabsContent value="tab2">Tab 2 content</TabsContent>
</Tabs>
```

### Table

```tsx
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "@/components/ui/table";

<Table>
  <TableHeader>
    <TableRow>
      <TableHead>Name</TableHead>
      <TableHead>Email</TableHead>
      <TableHead>Score</TableHead>
    </TableRow>
  </TableHeader>
  <TableBody>
    {users.map((user) => (
      <TableRow key={user.id}>
        <TableCell>{user.name}</TableCell>
        <TableCell>{user.email}</TableCell>
        <TableCell>{user.score}</TableCell>
      </TableRow>
    ))}
  </TableBody>
</Table>
```

### Toast (with Sonner)

```tsx
import { toast } from "sonner";

// Simple messages
toast.success("Success message");
toast.error("Error message");
toast.info("Info message");
toast.warning("Warning message");

// With description
toast.success("Quiz Completed", {
  description: "You scored 85%!",
});

// With action
toast("Quiz available", {
  action: {
    label: "Take Quiz",
    onClick: () => router.push("/quiz"),
  },
});
```

### Form with Validation

```tsx
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";

<form onSubmit={handleSubmit}>
  <div className="space-y-4">
    <div className="space-y-2">
      <Label htmlFor="email">Email</Label>
      <Input
        id="email"
        type="email"
        placeholder="user@example.com"
        {...register("email")}
      />
      {errors.email && (
        <p className="text-sm text-destructive">{errors.email.message}</p>
      )}
    </div>
    <Button type="submit">Submit</Button>
  </div>
</form>
```

## Customization

### Extending Variants

Add new variants using CVA (class-variance-authority):

```tsx
// button.tsx
const buttonVariants = cva(
  "inline-flex items-center justify-center ...",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground ...",
        // Add custom variant
        success: "bg-green-600 text-white hover:bg-green-700",
      },
    },
  }
);
```

### Custom Styling

Override with className:

```tsx
<Card className="border-primary/50 shadow-lg">
  <Button className="w-full bg-gradient-to-r from-blue-500 to-purple-500">
    Custom Button
  </Button>
</Card>
```

## Adding New Shadcn Components

Use the Shadcn CLI:

```bash
npx shadcn-ui@latest add [component-name]

# Examples
npx shadcn-ui@latest add calendar
npx shadcn-ui@latest add command
npx shadcn-ui@latest add tooltip
```

This will:
1. Add the component to `components/ui/`
2. Install required Radix UI dependencies
3. Update `components.json` if needed

## Related Documentation

- [Parent: Components Overview](../README.md)
- [Common Components](../common/README.md) - Components using these primitives
- [Shadcn/ui Docs](https://ui.shadcn.com/) - Full component documentation
- [Radix UI Docs](https://www.radix-ui.com/) - Underlying primitives
