import { HStack } from '@/components/HStack';
import { Text } from '@/components/Text';
import { VStack } from '@/components/VStack';
import { useOnScreenFocusCallback } from '@/hooks/useOnScreenFocusCallback';
import { ticketService } from '@/services/tickets';
import { TicketListData } from '@/types/ticket';
import { router, useNavigation } from 'expo-router';
import { useCallback, useEffect, useState } from 'react';
import { Alert, FlatList, TouchableOpacity } from 'react-native';

export default function TicketsScreen() {
  const navigation = useNavigation();

  const [isLoading, setIsLoading] = useState(false);
  const [tickets, setTickets] = useState<TicketListData>();

  function onGoToTicketPage(id: number) {
    router.push(`/(authed)/(tabs)/(tickets)/ticket/${id}`);
  }

  const fetchTickets = useCallback(async () => {
    try {
      setIsLoading(true);
      const response = await ticketService.getAll();
      setTickets(response.data[0]);
      console.log("hihi:",response.data[0])
      console.log("event:",response.data[0].tickets)
    } catch (error) {
      Alert.alert("Error", "Failed to fetch tickets");
    } finally {
      setIsLoading(false);
    }
  }, []);

  useOnScreenFocusCallback(fetchTickets);

  useEffect(() => {
    navigation.setOptions({ headerTitle: "Tickets" });
  }, [navigation]);

  return (
    <VStack flex={ 1 } p={ 20 } pb={ 0 } gap={ 20 }>

      <HStack alignItems='center' justifyContent='space-between'>
        <Text fontSize={ 18 } bold>{ tickets?.tickets?.length ?? 0 } Tickets</Text>
      </HStack>

      <FlatList
        keyExtractor={ (item) => item.ticket_id.toString() }
        data={ tickets?.tickets }
        onRefresh={ fetchTickets }
        refreshing={ isLoading }
        renderItem={ ({ item: ticket }) => (
          <TouchableOpacity disabled={ ticket.entered ?? false} onPress={ () => onGoToTicketPage(ticket.ticket_id) }>
            <VStack
              key={ ticket.ticket_id }
              {...({style:{
                height: 120,
                gap: 20,
                opacity: ticket.entered ? 0.5 : 1,
              }} as any)}
            >

              <HStack>
                <VStack
                  {...({style: {
                    height: 120,
                    width: "70%",
                    padding: 20,
                    justifyContent: 'space-between',
                    backgroundColor: "white",
                    borderTopLeftRadius: 20,
                    borderBottomLeftRadius: 20,
                    borderTopRightRadius: 5,
                    borderBottomRightRadius: 5,
                  }} as any )}
                >
                  <HStack alignItems='center'>
                    <Text fontSize={ 22 } bold >{ ticket.event.title }</Text>
                    <Text fontSize={ 22 } bold > | </Text>
                    <Text fontSize={ 16 } bold >{ ticket.event.location }</Text>
                  </HStack>
                  <Text fontSize={ 12 } >{ new Date(ticket.event.date).toLocaleString() }</Text>
                </VStack>

                <VStack
                  {...({style:{
                    height: 110,
                    width: "1%",
                    alignSelf: 'center',
                    borderColor: 'lightgray',
                    borderWidth: 2,
                    borderStyle: 'dotted',
                  }} as any)}
                />

                <VStack
                  {...({style:{
                    width: "25%",
                    justifyContent: 'center',
                    alignItems: 'center',
                    backgroundColor: "white",
                    borderTopRightRadius: 20,
                    borderBottomRightRadius: 20,
                    borderTopLeftRadius: 5,
                    borderBottomLeftRadius: 5,
                  }} as any)}
                >
                  <Text fontSize={ 16 } bold> { ticket.entered ? "Used" : "Available" } </Text>
                  {/* { ticket.entered && <Text mt={ 10 } fontSize={ 10 }>{ new Date(ticket.updatedAt).toLocaleString() }</Text> } */}
                </VStack>
              </HStack>

            </VStack>
          </TouchableOpacity>
        ) }
        ItemSeparatorComponent={ () => <VStack h={ 20 } /> }
      />

    </VStack>
  );
}