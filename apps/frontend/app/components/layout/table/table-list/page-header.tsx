import { Link } from "react-router";
import { Button } from "~/components/ui/button";
import { Plus } from "lucide-react";

export interface PageHeaderProps {
  title: string;
  description?: string;
  addButtonText?: string;
  addButtonLink?: string;
  onAddClick?: () => void;
}

export function PageHeader({
  title,
  description,
  addButtonText = "Add New",
  addButtonLink,
  onAddClick,
}: PageHeaderProps) {
  const renderButton = () => {
    const buttonContent = (
      <Button onClick={onAddClick}>
        <Plus className="mr-2 size-4" />
        {addButtonText}
      </Button>
    );

    if (addButtonLink) {
      return <Link to={addButtonLink}>{buttonContent}</Link>;
    }

    return buttonContent;
  };

  return (
    <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">{title}</h1>
        {description && (
          <p className="text-muted-foreground mt-2">{description}</p>
        )}
      </div>
      {(addButtonLink || onAddClick) && renderButton()}
    </div>
  );
}
