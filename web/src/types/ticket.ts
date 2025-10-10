export interface Ticket {
  id: string;
  title: string;
  description: string;
  mode: boolean;
  registered_count: number;
  quota: number;
  image: string;
}

export interface TicketMutation {
  title: string;
  description: string;
  mode: boolean;
  quota: number;
  image: string;
}
