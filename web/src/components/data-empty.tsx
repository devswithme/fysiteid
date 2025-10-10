import { cn } from "@/lib/utils";

interface DataEmptyProps {
  className?: string;
}

const DataEmpty = ({ className }: DataEmptyProps) => {
  return (
    <div
      className={cn(
        className,
        "w-full flex flex-col justify-center items-center"
      )}
    >
      <img src="/notfound.svg" />
      <h1 className="text-2xl text-center">
        Nothing to see here yet. Add some magic!
      </h1>
    </div>
  );
};

export default DataEmpty;
