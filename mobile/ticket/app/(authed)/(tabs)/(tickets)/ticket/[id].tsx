import { Text } from '@/components/Text';
import { VStack } from '@/components/VStack';
import { useOnScreenFocusCallback } from '@/hooks/useOnScreenFocusCallback';
import { ticketService } from '@/services/tickets';
import { Ticket } from '@/types/ticket';
import { router, useLocalSearchParams, useNavigation } from 'expo-router';
import { useCallback, useEffect, useState } from 'react';
import { Image } from 'react-native';

export default function TicketDetailsScreen() {
  const navigation = useNavigation();
  const { id } = useLocalSearchParams();

  const [ticketData, setTicketData] = useState<Ticket | null>(null);
  const [qrcode, setQrcode] = useState<string | null>(null);

  const fetchTicket = useCallback(async () => {
    try {
      const response = await ticketService.getOne(Number(id));
      setTicketData(response.data[0].ticket);
      setQrcode(response.data[0].qrcode);
    } catch (error) {
      router.back();
    }
  }, [id, router]);

  useOnScreenFocusCallback(fetchTicket);

  useEffect(() => {
    fetchTicket();
  }, [fetchTicket]);

  useEffect(() => {
    navigation.setOptions({ headerTitle: "" });
  }, [navigation]);

  if (!ticketData) return null;

  return (
    <VStack
      {...({style:{
        flex: 1,
        alignItems: 'center',
        gap: 20,
        margin: 20,
        padding: 20,
      }} as any)}
    >
      <Text fontSize={ 50 } bold >{ ticketData.event.title }</Text>
      <Text fontSize={ 20 } bold >{ ticketData.event.location }</Text>
      <Text fontSize={ 16 } color='gray'>{ new Date(ticketData.event.date).toLocaleString() }</Text>

      <Image
        style={ { borderRadius: 20 } }
        width={ 380 }
        height={ 380 }
        source={ { uri: `data:image/png;base64,${qrcode}` } }
      />
    </VStack>
  );
}