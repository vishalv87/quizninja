# Categories Components

## Overview

Components for browsing and displaying quiz categories. Users can explore quizzes organized by topic areas like Science, History, Technology, etc.

## Components

| Component | File | Purpose |
|-----------|------|---------|
| CategoryCard | `CategoryCard.tsx` | Single category card |
| CategoryGrid | `CategoryGrid.tsx` | Grid of category cards |

## CategoryCard

Card displaying a quiz category with actions.

### Props

| Prop | Type | Description |
|------|------|-------------|
| `category` | `Category` | Category data |
| `isFavorite` | `boolean` | Favorite status |
| `onToggleFavorite` | `(id) => void` | Toggle favorite |

### Display

- Category icon (image or fallback BookOpen icon)
- Display name
- Quiz count badge
- Description (2 lines max)
- Favorite heart button
- Browse and Quick Start buttons

### Visual Design

```tsx
// Gradient border wrapper
<div className="group p-[2px] rounded-2xl bg-gradient-to-br from-violet-500 to-indigo-500">
  <Card className="bg-white dark:bg-background rounded-[14px]">
    {/* Category content */}
  </Card>
</div>
```

### Actions

```tsx
// Browse - navigate to category page
<Button asChild>
  <Link href={`/quizzes/category/${category.id}`}>
    Browse
    <ArrowRight className="ml-2 h-4 w-4" />
  </Link>
</Button>

// Quick Start - random quiz from category
<Button onClick={handleQuickStartClick}>
  <Shuffle className="mr-2 h-4 w-4" />
  Quick Start
</Button>
```

### Image Handling

```tsx
// With fallback for missing/error images
{category.icon_url && !imageError ? (
  <Image
    src={category.icon_url}
    alt={category.display_name}
    onError={() => setImageError(true)}
  />
) : (
  <BookOpen className="h-6 w-6 text-violet-700" />
)}
```

### Usage

```tsx
import { CategoryCard } from "@/components/categories/CategoryCard";

<CategoryCard
  category={category}
  isFavorite={favoriteCategories.includes(category.id)}
  onToggleFavorite={handleToggleFavorite}
/>
```

---

## CategoryGrid

Grid layout for displaying categories.

### Props

| Prop | Type | Description |
|------|------|-------------|
| `categories` | `Category[]` | Categories array |
| `favorites` | `string[]` | Favorite category IDs |
| `onToggleFavorite` | `(id) => void` | Toggle favorite |
| `isLoading` | `boolean` | Loading state |

### Features

- Responsive grid (1-4 columns)
- Loading skeleton state
- Empty state message

### Grid Layout

```tsx
<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
  {categories.map(category => (
    <CategoryCard
      key={category.id}
      category={category}
      isFavorite={favorites.includes(category.id)}
      onToggleFavorite={onToggleFavorite}
    />
  ))}
</div>
```

### Usage

```tsx
import { CategoryGrid } from "@/components/categories/CategoryGrid";
import { useCategories } from "@/hooks/useCategories";

function CategoriesPage() {
  const { data: categories, isLoading } = useCategories();
  const [favorites, setFavorites] = useState<string[]>([]);

  return (
    <CategoryGrid
      categories={categories}
      favorites={favorites}
      onToggleFavorite={(id) => toggleFavorite(id)}
      isLoading={isLoading}
    />
  );
}
```

---

## Categories Page Structure

```tsx
// /categories page
<div className="space-y-8">
  {/* Page Header */}
  <div>
    <h1>Quiz Categories</h1>
    <p>Browse quizzes by topic</p>
  </div>

  {/* Search/Filter */}
  <div className="flex gap-4">
    <Input placeholder="Search categories..." />
    <Select>
      <SelectItem value="all">All Categories</SelectItem>
      <SelectItem value="favorites">Favorites Only</SelectItem>
    </Select>
  </div>

  {/* Category Grid */}
  <CategoryGrid
    categories={filteredCategories}
    favorites={userFavorites}
    onToggleFavorite={handleToggleFavorite}
    isLoading={isLoading}
  />
</div>
```

## Data Types

```typescript
interface Category {
  id: string;
  name: string;
  display_name: string;
  description?: string;
  icon_url?: string;
  quiz_count: number;
  created_at: string;
}

interface CategoryFilters {
  search?: string;
  favoritesOnly?: boolean;
}
```

## Hooks Used

```typescript
// Get all categories
const { data: categories, isLoading } = useCategories();

// Get single category with quizzes
const { data: category } = useCategory(categoryId);

// Get quizzes by category
const { data: quizzes } = useQuizzes({ category: categoryName });
```

## Route Structure

| Route | Purpose |
|-------|---------|
| `/categories` | All categories grid |
| `/quizzes/category/[categoryId]` | Quizzes in specific category |
| `/quizzes?category=[name]` | Filtered quiz list |

## Related Documentation

- [Parent: Components Overview](../README.md)
- [Quiz Components](../quiz/README.md)
- [Category Types](../../types/README.md)
- [useCategories Hook](../../hooks/README.md)

