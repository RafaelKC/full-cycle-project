import { Controller, Get, Query } from '@nestjs/common';
import { PlacesService } from './places.service';

@Controller('places')
export class PlacesController {
  constructor(private readonly placeService: PlacesService) {}

  @Get()
  public async getPlaces(@Query('text') text: string) {
    return await this.placeService.getPlaces(text);
  }
}
