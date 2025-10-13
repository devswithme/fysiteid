// import TicketForm from "@/components/form/ticket";
import Modal from "@/components/modal";
import { Badge } from "@/components/ui/badge";
import { Button, buttonVariants } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  useGenerateState,
  useGetRegistrantByTicketID,
} from "@/hooks/use-registrant";
import { useTicketDetail } from "@/hooks/use-ticket";
import { downloadCSV, formatDate } from "@/lib/utils";
import { ticketService } from "@/service/ticket";
import {
  createFileRoute,
  Link,
  notFound,
  useNavigate,
} from "@tanstack/react-router";
import { ArrowLeft, QrCode, Save } from "lucide-react";
import { QRCodeSVG } from "qrcode.react";
import { useEffect, useState } from "react";
import { toast } from "sonner";

function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedValue(value), delay);
    return () => clearTimeout(timer);
  }, [value, delay]);

  return debouncedValue;
}

export const Route = createFileRoute("/_auth/dash/ticket/$id/")({
  validateSearch: (search: Record<string, unknown>) => {
    return {
      page: Number(search.page) || 1,
    };
  },
  loader: async ({ params: { id } }) => {
    try {
      const data = await ticketService.getTicketById(id);
      if (!data) throw notFound();

      return data;
    } catch {
      throw notFound();
    }
  },
  component: RouteComponent,
});

function RouteComponent() {
  const [searchTerm, setSearchTerm] = useState("");
  const debouncedSearch = useDebounce(searchTerm, 500);

  const navigate = useNavigate();
  const { page } = Route.useSearch();

  const [qrOpen, setQrOpen] = useState(false);

  const { id } = Route.useParams();
  const { data: ticket } = useTicketDetail(id);

  const { data: registrants } = useGetRegistrantByTicketID({
    ticketID: id,
    page,
    search: debouncedSearch,
  });

  const { refetch: fetchAllData } = useGetRegistrantByTicketID({
    ticketID: id,
    page: 0,
    limit: 0,
    search: debouncedSearch,
  });

  const { data: state, refetch: regenerate } = useGenerateState(id);

  const handlePageChange = (newPage: number) => {
    navigate({
      to: "/dash/ticket/$id",
      params: { id },
      search: {
        page: newPage,
      },
    });
  };

  const handleExportCSV = async () => {
    try {
      const { data } = await fetchAllData();
      if (data?.data) {
        downloadCSV(
          data.data,
          `${ticket?.title}-${new Date().toISOString()}.csv`
        );
      }
    } catch {
      toast.error("Failed to export data");
    }
  };

  return (
    <>
      <div className="space-y-6">
        <div className="flex justify-between">
          <Link
            to="/dash/ticket"
            className={buttonVariants({ variant: "secondary" })}
          >
            <ArrowLeft /> Back
          </Link>
          <Button onClick={handleExportCSV}>
            <Save /> Export CSV
          </Button>
        </div>
        <div>
          <h1 className="text-2xl font-semibold">{ticket?.title}</h1>
          <Badge variant="secondary">
            {ticket?.registered_count}/{ticket?.quota}
          </Badge>
        </div>
        <div className="flex gap-2 flex-wrap">
          {ticket?.mode && (
            <Modal
              trigger={
                <Button
                  onClick={async () => {
                    await regenerate();
                    setQrOpen(true);
                  }}
                >
                  <QrCode /> Generate QR
                </Button>
              }
              open={qrOpen}
              onOpenChange={setQrOpen}
            >
              <div className="flex justify-center items-center">
                <QRCodeSVG
                  value={`${
                    import.meta.env.VITE_APP_URL
                  }/ticket/${id}?state=${state}`}
                  className="w-fit h-fit p-4 bg-muted border aspect-square"
                />
              </div>
            </Modal>
          )}
        </div>
      </div>
      <div className="space-y-4">
        <Input
          placeholder="Search registrants.."
          className="max-w-sm"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Registrant</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Time</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {!registrants?.data.length && (
              <TableCell colSpan={3} className="p-4 text-center">
                No data
              </TableCell>
            )}
            {registrants?.data?.map((registrant) => (
              <TableRow key={registrant.username}>
                <TableCell className="font-medium flex items-center gap-x-2">
                  <img
                    src={registrant.picture}
                    className="size-5 rounded-full"
                  />
                  {registrant.username}
                </TableCell>
                <TableCell>
                  <Badge
                    variant={registrant.is_verified ? "success" : "destructive"}
                  >
                    {registrant.is_verified ? "Verified" : "Unverified"}
                  </Badge>
                </TableCell>
                <TableCell className="text-right">
                  {formatDate(registrant.created_at)}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <Pagination>
          <PaginationContent>
            {page >= 1 && (
              <PaginationItem>
                <PaginationPrevious
                  href="#"
                  onClick={() => handlePageChange(Math.max(1, page - 1))}
                />
              </PaginationItem>
            )}
            {Array.from({ length: registrants?.total_pages || 0 }, (_, i) => (
              <PaginationItem key={i + 1}>
                <PaginationLink
                  href="#"
                  onClick={() => handlePageChange(i + 1)}
                  isActive={page === i + 1}
                >
                  {i + 1}
                </PaginationLink>
              </PaginationItem>
            ))}
            {page <= (registrants?.total_pages || 1) && (
              <PaginationItem>
                <PaginationNext
                  href="#"
                  onClick={() =>
                    handlePageChange(
                      Math.min(registrants?.total_pages || 1, page + 1)
                    )
                  }
                />
              </PaginationItem>
            )}
          </PaginationContent>
        </Pagination>
      </div>
    </>
  );
}
